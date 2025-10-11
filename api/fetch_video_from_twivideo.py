import os
import requests
import json
import re
import time
import mysql.connector
from mysql.connector.errors import DatabaseError, InterfaceError
from bs4 import BeautifulSoup, Comment
from celery import Celery
from celery.schedules import crontab
import cloudscraper
import brotli
import ulid

def generate_ulid():
    return ulid.new()

URL = "https://twivideo.net/templates/view_lists.php"

# HTTPヘッダー（403回避用）
HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
    "Referer": "https://twivideo.net/?ranking",
    "X-Requested-With": "XMLHttpRequest",
    "Origin": "https://twivideo.net",
    "Accept": "*/*",
    "Accept-Encoding": "gzip, deflate, br",
    "Accept-Language": "ja,en-US;q=0.9,en;q=0.8",
    "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
    "Cookie": "_ga=GA1.1.1547131204.1737202865; _im_vid=01JHWNR4FC3FR55SKSBTNHQD3H; jp1_et_freq={\"4313\":[0,1746903714395]}; PHPSESSID=u6iuv254qkalpeg9gja0tpi0ei; __bnc_pfpuid__=5a-X05UlYIicf; cf_clearance=PuI0WUs_VwZpWjte6jyjxdscDZ4PFdeFgDnAYaKEYUc-1747572174-1.2.1.1-s129WN7UrH1Ysqia72fzdtzembKpSlOAGE._AAheu7LOMgcCBPefZi3jh.GzumtsJk1EuD7c3DO1pvEwhlRn4Ge7VefUsRfbf5DGBbV6oagPOgtwrLyK0tehUTP1pJlDgR6q1uHcs.vhGMqdKch7FXks7Rf3xJVrd1XDoOKBEFrN00k_GZz7uhHASRJIyahLwGB9euKQbq6ZJLz7cvSinYB.MqsC2rb3QBN3HLGihnCUbvV0RSVOV2DDB9Tc8FRIz5UVHo96UebQF.JlDunmhGqeoN4TUUOsIeguJyKAbkXMyaVkEJIgI2Eu8kkRR0J71E4fXDfjLAPeLNKqEtwTg9TvjYE9kO5iespsWDTCr.I; _ga_RG2YW44R3P=GS2.1.s1747572167$o30$g1$t1747572856$j0$l0$h0"
}

# POSTリクエストのデータ
DATA = {
    "offset": 0,
    "limit": 1000,
    "tag": "null",
    "type": "ranking",
    "order": "24",
    "le": 1000,
    "ty": "p6",
    "myarray": "[]",
    "offset_int": 0
}
def fetch_video_list():
    scraper = cloudscraper.create_scraper()

    print("fetching video list...")
    response = scraper.post(URL, headers=HEADERS, data=DATA)
    
    print(f"{response.status_code}")
    if response.status_code != 200:
        print(f"{response.status_code}")
        return f"エラー: {response.status_code}"

    print("saving the response as text ")
    # レスポンスをtextとして保存
    response.encoding = response.apparent_encoding
    with open("response.html", "w", encoding="utf-8") as f:
        f.write(response.text)

    soup = BeautifulSoup(response.text, "html.parser")
    videos = []
    # 動画URLからIDを取り出すための正規表現
    pattern = r"/(?:ext_tw_video_thumb|amplify_video_thumb)/(\d+)/"

    for item in soup.find_all("div", class_="art_li"):
        video_data = {}
        
        # 動画のURL
        twivideo_url = item.find("a", class_="item_clk item_link")
        # headers = {
        #     "User-Agent": "Mozilla/5.0",
        #     "Referer": "https://twivideo.net/"
        # }
        print(twivideo_url.get("href"))
        video_data["video_url"] = twivideo_url.get("href")
        # response = requests.get(twivideo_url.get("href"), headers=headers, allow_redirects=False)
        # print(f"URL: {twivideo_url.get('href')}")  # デバッグ用
        # print(f"CODE: { response.status_code}")  # デバッグ用
        # if response.status_code in [301, 302, 200]:
        #     mp4_url = response.headers.get("Location")
        #     if mp4_url:
        #         print(f"MP4 URL: {mp4_url}")  # デバッグ用
        #         video_data["video_url"] = mp4_url
        #     else :
        #         print("❌ リンク取得失敗。Locationヘッダーが見つかりません。")
                
        # else:
        #     print(f"❌ リンク取得失敗。ステータスコード: {response.status_code}")
        
        
            
        # 動画URLから動画IDを取得
        # match = re.search(pattern, video_data["video_url"])
        # if match:
        #     video_data["id"] = match.group(1)
        # else:
        #     video_data["id"] = str(generate_ulid())
            
        # サムネイルのURL
        thumbnail = item.find("img")
        if thumbnail:
            video_data["thumbnail_url"] = thumbnail.get("src")

        match = re.search(pattern, video_data["thumbnail_url"])
        if match:
            video_data["id"] = match.group(1)
        else:
            video_data["id"] = str(generate_ulid())

        # ランキング（"No.1" → "1" に変換）
        ranking = item.find("div", class_="item_ranking")
        if ranking:
            try:
                video_data["ranking"] = int(ranking.text.strip().replace("No.", "").strip())
            except Exception:
                video_data["ranking"] = 0
        # DL数（コメント内の数字を取得）
        download_count = 0
        for comment in item.find_all(string=lambda text: isinstance(text, Comment)):
            m = re.search(r'(\d+)DL', comment)
            if m:
                download_count = int(m.group(1))
                break
        video_data["download_count"] = download_count

        # Twitter投稿ページのURL（tweet_urlとして格納）
        tw_icon_div = item.find("div", class_="tw_icon")
        tweet_url = None
        if tw_icon_div:
            a_tag = tw_icon_div.find("a", href=True, rel="noopener noreferrer")
            if a_tag and "x.com" in a_tag["href"]:
                tweet_url = a_tag["href"]
        video_data["tweet_url"] = tweet_url

        videos.append(video_data)

    # 重複するIDがあれば、最初のものだけを残す
    unique_videos = {}
    deduped_videos = []
    for video in videos:
        vid = video.get("id")
        if vid not in unique_videos:
            unique_videos[vid] = True
            deduped_videos.append(video)
    videos = deduped_videos
    
    # JSONに保存（デバッグ用）
    with open("videos.json", "w", encoding="utf-8") as json_file:
        json.dump({ "videos": videos }, json_file, indent=4, ensure_ascii=False)
    
    #videosをjsonにしてAPIにPOST
    with open("videos.json", "r", encoding="utf-8") as json_file:
        videos_json = json.load(json_file)
        
    # videos_jsonをAPIにPOST
    url = "http://localhost:8000/api/videos/create-bulk"
    
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, headers=headers, json=videos_json)
    
    if response.status_code == 200:
        print("動画リストをAPIにPOSTしました。")
    else:
        print(f"APIへのPOSTに失敗しました。ステータスコード: {response.status_code}")
        
    

if __name__ == "__main__":
    fetch_video_list()