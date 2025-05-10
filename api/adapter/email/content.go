package email

import (
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"
)

func ContentToResetPassword(token string) (subject, body string) {
	subject = "【e-privado】パスワード再設定のお知らせ"
	body = `
e-privadoをご利用いただきありがとうございます。
下記URLにアクセスし、パスワードを再設定を完了させてください。

%s

※24時間以内に手続きを完了しない場合、上記URLは無効になります。
最初から手続きをやり直してください。
`
	body = fmt.Sprintf(body, config.FrontendURL()+"/password-change?token="+token)
	return
}

func ContentToUpdateEmail(token string) (subject, body string) {
	subject = "【e-privado】メールアドレス変更受付のお知らせ"
	body = `
e-privadoをご利用いただきありがとうございます。
下記URLにアクセスし、メールアドレスの変更を完了させてください。

%s

※24時間以内に手続きを完了しない場合、上記URLは無効になります。
最初から手続きをやり直してください。
`
	body = fmt.Sprintf(body, config.FrontendURL()+"/email/update/complete?token="+token)

	return
}

// パッケージを購入したが、定員が満員のため、キャンセルされた場合のメール
func ContentToPackagePlanPaymentCancelForUser(userName, planName, providerPriceID string, price int32) (subject, body string) {
	subject = "【重要】パッケージプラン購入に関するお知らせ"
	body = `
	%s 様

	この度は【e-privado】のパッケージプランを購入していただき、誠にありがとうございます。
	しかしながら、大変申し訳ございませんが、お申込みいただいたプランの定員が既に上限に達しており、
	お客様の購入をお受けすることができませんでした。

	【プラン内容】
	・請求ID: %s
	・プラン名: %s
	・価格: %d円（税込）

	お客様にはご迷惑をおかけし、心よりお詫び申し上げます。
	お支払いいただいた料金は、全額、速やかにご指定の口座へ返金致します。
	返金の手続きには数日を要する場合がございますが、何卒ご了承ください。

	引き続き【e-privado】をご利用いただけますよう、よろしくお願いいたします。
	`

	body = fmt.Sprintf(body, userName, providerPriceID, planName, price)
	return
}

// パッケージプランの購入に失敗したため、ユーザーに返金する必要があることを管理者に通知するメール
func ContentToPackagePlanPaymentCancelForAdmin(userName, planName, providerPriceID string, price int32) (subject, body string) {
	subject = "【重要】パッケージプランの返金対応が必要です。"
	body = `
	e-privadoシステム管理者様

	下記ユーザーのパッケージプランの購入が定員超過によりキャンセルされました。
	返金の手続きをお願いいたします。

	【プラン内容】
	・購入ユーザー: %s 様
	・請求ID: %s
	・プラン名: %s
	・価格: %d円（税込）

	何卒、よろしくお願いいたします。
	`
	body = fmt.Sprintf(body, userName, providerPriceID, planName, price)
	return
}

func ContentToExportCsvComplete(presignedUrl string) (subject, body string) {
	subject = "【e-privado】CSVエクスポート完了のお知らせ"
	body = `
e-privadoをご利用いただきありがとうございます。
下記URLにアクセスすることで、csvをダウンロードできます。

%s

※24時間で上記URLは無効になります。
再びアクセスするには、最初から手続きをやり直してください。
`
	body = fmt.Sprintf(body, presignedUrl)
	return
}
