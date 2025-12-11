import os
import json
import re
import time
try:
    import cloudscraper
except ImportError:
    cloudscraper = None
import requests
from bs4 import BeautifulSoup, Comment
import ulid
from dotenv import load_dotenv

load_dotenv()

# =====================================================
# 設定
# =====================================================

# ローカルに保存したランキングHTML
TWIVIDEO_HTML_PATH = "response.html"

# Cookie を複数本持っておく（ここにさっきの2本を入れる）
TWIVIDEO_COOKIES = [
    # 1本目
    """_ga=GA1.1.1547131204.1737202865; _im_vid=01JHWNR4FC3FR55SKSBTNHQD3H; jp1_et_freq={"4313":[0,1746903714395]}; uid=diu5bf1012b4a0912f4a87e96fcc866ecda; __bnc_pfpuid__=5a-0LhV4QFPXU; is_valid_param=1d63724ec871ee538a26b920604e857c; PHPSESSID=8j3245pgo56tpavfehk6f1mkl3; cf_clearance=FFu6kqsThVPd1xr4gnRXZUgWz3i0YLdpNhHPVdeyJ2s-1765146241-1.2.1.1-zjVV01jdMlnyEFULS5WNoBxE9bYyKXInDxUzxDwHVoS24p91nnb_2AC_tr3f0sLtrtVarlh.tFfHu757QVt8HgIV4PgeNpZ0u0XpNCaKlkHul4M3WL2oFxBydOsqIfQXuYRWQ0ysveZ2KNHItgiXMgs8XUMm8xcnNXyNJk7Exw.Fg4g2an3PFVCvNyseKOvMOaqs7J4879xNvoG8HYQW8.baeYEZ6VoZZs90aBJMgPc; _ga_RG2YW44R3P=GS2.1.s1765146240$o71$g1$t1765146271$j29$l0$h0""",
    # 2本目（今送ってくれたやつ）
    """_ga=GA1.1.1935914182.1764180130; _im_vid=01KB0N8YERYTBBTM90WDTFADMK; is_valid_param=dad43fcc1bb8cf95a3b01a04174b2185; PHPSESSID=np2vec3clhhpgv3l94fo5r607m; __bnc_pfpuid__=5a-IrCV4jjvfF; cf_clearance=C_2dRqvVGojxEMkrJtPbRvNSwVTmyjk5yCCqkTTk_DQ-1765147716-1.2.1.1-A.WVTkx49q1EWWAtg1DVDRwRG0YvmpPtlhHU37EtHqm9defS22rDdxQugmDcO5WedH1Z7prvFKCxcaJvoVMKQzqytat2He7u.ag30FkpZovRGG5xxeCRIp2WtLbr5mKQxIm5CePxHyRZCPjcMQVseYMJ72keBVIV8YCq7hGG_Zm8sF4Py53DhM.yIYgDdaWIlAEKhX2C5KIFc0jBLtOpK1H961GygNxuNAgiVF_8VdA; _ga_RG2YW44R3P=GS2.1.s1765147716$o2$g1$t1765147722$j54$l0$h0""",
]

TWIVIDEO_LINK_API_BASE = "https://twivideo.net/api/link.php"

# link.php 叩くときのヘッダー（Cookie は後で差し込む）
LINK_HEADERS_BASE = {
    "User-Agent": (
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
        "AppleWebKit/537.36 (KHTML, like Gecko) "
        "Chrome/123.0.0.0 Safari/537.36"
    ),
    "Referer": "https://twivideo.net/?ranking",
    "Accept": "*/*",
    "Accept-Language": "ja,en-US;q=0.9,en;q=0.8",
}

# Xroll の API
DEFAULT_API_URL = "https://xroll.net/api/videos/create-bulk"
XROLL_API_ENV = "XROLL_API_URL"

# 429 回避用
LINK_MAX_CALLS = 1000
LINK_INTERVAL_SEC = 0.7         # 1件ごとのスリープ
LINK_RETRY_BASE_WAIT = 5.0      # 429 のときの待ち秒
LINK_MAX_RETRIES = 3

# サムネイルURLから Twitter video id を抜くパターン（元コードのまま）
THUMB_ID_PATTERN = r"/(?:ext_tw_video_thumb|amplify_video_thumb)/(\d+)/"


def generate_ulid_str() -> str:
    return str(ulid.new())


# =====================================================
# HTML ロード
# =====================================================

def load_html_from_file() -> str:
    path = TWIVIDEO_HTML_PATH.strip()
    if not path:
        raise RuntimeError(
            f"{TWIVIDEO_HTML_PATH} が設定されていません。"
            " .env に TWIVIDEO_HTML_PATH=ダウンロードしたHTMLファイルのパス を書いて。"
        )

    if not os.path.exists(path):
        raise FileNotFoundError(f"TWIVIDEO_HTML_PATH で指定されたファイルが存在しません: {path}")

    print(f"[HTML] load from file: {path}")
    with open(path, "r", encoding="utf-8") as f:
        return f.read()


# =====================================================
# link.php から tweet_url を取る処理
# =====================================================

def build_link_headers(cookie: str) -> dict:
    headers = dict(LINK_HEADERS_BASE)
    cookie = (cookie or "").strip()
    if not cookie:
        print("[WARN] Cookie が空です。link.php で 403/429 の可能性あり。")
    else:
        headers["Cookie"] = cookie
    return headers


def fetch_tweet_url_from_id(session: requests.Session, video_id: str, headers: dict) -> str | None:
    """
    twivideo の link.php から tweetURL を取得
    （TwivideoSource._fetch_tweet_url_from_id をほぼそのまま移植）
    """
    api_url = f"{TWIVIDEO_LINK_API_BASE}?id={video_id}"

    for attempt in range(1, LINK_MAX_RETRIES + 1):
        try:
            print(
                f"[LINK] link.php access id={video_id} "
                f"(attempt {attempt}/{LINK_MAX_RETRIES})"
            )
            resp = session.get(api_url, headers=headers, timeout=10)

            # 429 → 少し待って再試行
            if resp.status_code == 429:
                wait_sec = LINK_RETRY_BASE_WAIT * attempt
                print(f"[LINK][WARN] 429 Too Many Requests → sleep {wait_sec}s")
                time.sleep(wait_sec)
                continue

            if resp.status_code == 403:
                print(f"[LINK][WARN] 403 Forbidden for id={video_id}")
                return None

            resp.raise_for_status()

            text = resp.text.strip()
            if not text:
                print(f"[LINK][WARN] empty response from link.php id={video_id}")
                return None

            first_line = text.splitlines()[0].strip()
            if first_line.startswith("http") and (
                "x.com" in first_line or "twitter.com" in first_line
            ):
                print(f"[LINK] tweet_url via link.php OK: {first_line}")
                return first_line

            print(f"[LINK][WARN] unexpected link.php response: {first_line!r}")
            return None

        except requests.HTTPError as e:
            print(f"[LINK][WARN] HTTPError link.php id={video_id}: {e}")
            return None
        except Exception as e:
            print(f"[LINK][WARN] failed link.php id={video_id}: {e}")
            if attempt == LINK_MAX_RETRIES:
                return None
            sleep_sec = 2.0 * attempt
            print(f"[LINK] retry after {sleep_sec}s")
            time.sleep(sleep_sec)

    return None

def fetch_video_list_from_local_html():
    # ----- HTML 読み込み -----
    html = load_html_from_file()
    soup = BeautifulSoup(html, "html.parser")

    # セッション作成（cloudscraper 優先）
    if cloudscraper is not None:
        print("[SESSION] use cloudscraper")
        session = cloudscraper.create_scraper()
    else:
        print("[SESSION] use requests.Session (cloudscraper無し)")
        session = requests.Session()

    pattern = re.compile(THUMB_ID_PATTERN)

    videos = []
    link_calls = 0

    print("[PARSE] find <div class='art_li'> ...")
    # ★ ここで enumerate 追加（何番目か追えるように）
    for idx, item in enumerate(soup.find_all("div", class_="art_li"), start=1):
        print("\n==============================")
        print(f"[ITEM] index={idx}")

        video_data = {}

        # -------------------------
        # Twivideo 内部 ID（link.php 用）
        # -------------------------
        link_tag = item.find("a", class_="item_clk item_link")
        twivideo_id = None
        if link_tag:
            twivideo_id = link_tag.get("data-id")  # TwivideoSource と同じ
            # 一応 Twivideo のページURLも残しておく
            video_data["video_url"] = link_tag.get("href")
        else:
            video_data["video_url"] = None

        # -------------------------
        # サムネイル URL & そこから id を作る（元コード準拠）
        # -------------------------
        thumbnail = item.find("img")
        if thumbnail:
            video_data["thumbnail_url"] = thumbnail.get("src")
        else:
            video_data["thumbnail_url"] = None

        vid_id = None
        if video_data["thumbnail_url"]:
            m = pattern.search(video_data["thumbnail_url"])
            if m:
                vid_id = m.group(1)

        # ext_tw_video_thumb から取れなければ ULID
        if not vid_id:
            vid_id = generate_ulid_str()
        video_data["id"] = vid_id

        # -------------------------
        # ランキング
        # -------------------------
        ranking = 0
        ranking_div = item.find("div", class_="item_ranking")
        if ranking_div and ranking_div.text:
            try:
                ranking = int(ranking_div.text.strip().replace("No.", "").strip())
            except Exception:
                ranking = 0
        video_data["ranking"] = ranking

        # -------------------------
        # DL 数（コメント内の「○DL」）
        # -------------------------
        download_count = 0
        for comment in item.find_all(string=lambda text: isinstance(text, Comment)):
            m = re.search(r"(\d+)DL", comment)
            if m:
                download_count = int(m.group(1))
                break
        video_data["download_count"] = download_count

        # -------------------------
        # tweet_url を link.php から取得
        # -------------------------
        tweet_url = None

        if twivideo_id and link_calls < LINK_MAX_CALLS and TWIVIDEO_COOKIES:
            # ここで Cookie ローテーション
            cookie_idx = link_calls % len(TWIVIDEO_COOKIES)
            cookie_str = TWIVIDEO_COOKIES[cookie_idx]
            link_headers = build_link_headers(cookie_str)

            print(f"[PARSE] call link.php id={twivideo_id} cookie_idx={cookie_idx}")
            tweet_url = fetch_tweet_url_from_id(session, twivideo_id, link_headers)
            link_calls += 1
            time.sleep(LINK_INTERVAL_SEC)

        elif twivideo_id and link_calls >= LINK_MAX_CALLS:
            print(
                f"[PARSE] link.php call limit reached ({LINK_MAX_CALLS}) "
                f"→ skip id={twivideo_id}"
            )

        # fallback: HTML から x.com を拾う（念のため）
        if not tweet_url:
            tw_icon_div = item.find("div", class_=lambda cls: cls and "tw_icon" in cls)
            search_scope = tw_icon_div if tw_icon_div is not None else item
            a_tag = search_scope.find(
                "a",
                href=lambda h: h and ("x.com" in h or "twitter.com" in h),
            )
            if a_tag:
                tweet_url = a_tag["href"]

        video_data["tweet_url"] = tweet_url

        # ★ ここで、その場でできあがった JSON を出す
        print("[ITEM] video_data =")
        print(json.dumps(video_data, ensure_ascii=False, indent=2))

        videos.append(video_data)

        # 50個ごとにPOST
        if len(videos) >= 50:
            # -------------------------
            # Xroll API へ POST
            # -------------------------
            api_url = os.environ.get(XROLL_API_ENV, DEFAULT_API_URL)
            print(f"[POST] endpoint: {api_url}")
            headers = {"Content-Type": "application/json"}

            resp = requests.post(api_url, headers=headers, json={"videos": videos}, timeout=60)
            print(f"[POST] status: {resp.status_code}")
            if resp.status_code == 200:
                print("[POST] 動画リストをAPIにPOSTしました。")
                # POST後、videosリストをクリア
                videos = []
            else:
                print(f"[POST][WARN] APIへのPOSTに失敗しました。status={resp.status_code}")
                print(resp.text[:1000])

    # 最後に残りの動画もPOST
    if len(videos) > 0:
        print("[POST] sending remaining videos...")
        api_url = os.environ.get(XROLL_API_ENV, DEFAULT_API_URL)
        headers = {"Content-Type": "application/json"}
        resp = requests.post(api_url, headers=headers, json={"videos": videos}, timeout=60)
        print(f"[POST] status: {resp.status_code}")
        if resp.status_code == 200:
            print("[POST] 最後の動画リストをAPIにPOSTしました。")
        else:
            print(f"[POST][WARN] APIへのPOSTに失敗しました。status={resp.status_code}")
            print(resp.text[:1000])

if __name__ == "__main__":
    fetch_video_list_from_local_html()
