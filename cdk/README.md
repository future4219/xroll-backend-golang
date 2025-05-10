# CDK TypeScript

## Useful commands

- `npm run build` TS ファイルを JS にコンパイルする
- `npm run watch` 変更を監視して自動でコンパイルする
- `npm run test` テストの実行
- `cdk deploy` デプロイする。設定がなければデフォルトのプロファイルで実行される。
- `cdk diff` ソースコードとデプロイされた環境の diff を表示
- `cdk synth` CloudFormation のテンプレートファイルを生成する
- `npm run lint` リントの実行
- `npm run format` フォーマットの実行

## 初期設定

1. cdk のインストール

   ```sh
   npm install -g aws-cdk`
   ```

2. ソースコードの更新
   [テンプレートリポジトリ](https://gitlab.com/digeon-inc/templates/cdk-ecs-rds-template) の cdk 下の内容を案件のリポジトリに反映する。
   以下に気をつける
   テンプレートと案件リポジトリでディレクトリ構成が異なる場合がある
   モニタリングバッチなど、テンプレートで用意されているが案件では使わないものがある

3. アラート受入用に slack チャンネルを作成する
   stg 環境と prod 環境はそれぞれ別チャンネルとして作成する
   stg についてはそもそもアラート対応が必要か PM と確認
   Digeon の workspace id と、作成したチャンネルの id を控えておく

4. `./cdk/cdk.json`の設定を行う

5. デプロイ先の AWS アカウントでコンソール経由で設定する

   - ホストゾーン ID の取得 or PM に作成依頼
   - [ACM](https://ap-northeast-1.console.aws.amazon.com/acm/home?region=ap-northeast-1#/certificates/request) で証明書を依頼
   - IAM からアクセスキーを生成(値を控えておく)

6. デプロイ先の AWS にアクセス

   ```sh
   cd ./cdk
   cdk bootstrap
   ```

7. Gitlab の設定で CICD の variable を追加する
   - AWS_ACCOUNT_ID
   - AWS_ACCESS_KEY_ID
   - AWS_SECRET_ACCESS_KEY
   - AWS_DEFAULT_REGION

## Usage

### シークレット登録のための初回手動デプロイ

1. cdk で Secrets Manager のスタックをデプロイする。
   1. `cd cdk`
   2. `npx cdk deploy <アプリ名>-stg-sm-stack -c tag=tmp -c stage=stg --profile $AWS_PROFILE`
2. Sectrets Manager の環境変数を設定する。
   1. AWS コンソールにアクセス
   2. `アプリ名`-stg-secret を選択
   3. 『シークレットの値』に環境変数を入力する。一括で入力する場合 json 形式で入力すると楽。
3. cdk でほかのスタックをデプロイする。
   1. `cd cdk`
   2. `npx cdk deploy --all -c tag=tmp -c stage=stg --profile $AWS_PROFILE`
4. 以後の更新は ci によるデプロイに任せる

#### (本番環境用)監視アラート用の Slack の設定

1. [AWS Chatbot](https://us-east-2.console.aws.amazon.com/chatbot/home?region=ap-northeast-1#/chat-clients)内に Digeon のワークスペースがなければ作成をする

   `＊チャンネルとの接続はCDKで行うのでコンソール上で作成する必要はない`

2. cdk.json に作成したチャンネルの情報(slackWorkspaceId と slackChannelId を記入する)

   デフォルトでは監視用アラートは本番環境のみ作成されるが、cdk.json の alertMonitoringEnabled を変更することで stg 環境もアラート監視を行うことができる

   ただし、stg 環境と prod 環境を同じ Slack チャンネルで監視することはできないので、それぞれ別でチャンネルを用意しておくこと

### ドメインの設定

1. ホストゾーン ID 及びホスト名を取得
   ホストゾーン ID は Z から始まる 21 文字の文字列である。
2. ACM で証明書の作成
   完全修飾ドメイン名は`*.your-domain.com`とする。
   証明書の ARN を控えておく。
3. `cdk.json`の ecs の設定以下を追加する。
   ```json
   {
     "zoneName": "your-domain.com",
     "hostedZoneId": "Z0123456789ABCDEFGHIJ",
     "apiDomainName": "api.your-domain.com",
     "certificateArn": "arn:aws:acm:ap-northeast-1:012345678901:certificate/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
   }
   ```

### マイグレーションの実行

DB スキーマ変更時のマイグレーションは完全自動化ではなく、cdk デプロイによって AWS 上に生成されたマイグレーション用のタスク定義を用いて[AWS コンソール上](https://ap-northeast-1.console.aws.amazon.com/ecs/v2/task-definitions?region=ap-northeast-1)から実行する

1. マイグレーション用のタスク定義(app-{stg,prod}-**task-migration-def-family**)を選択し、最新のリビジョンを選択する
2. コンソール右上の「デプロイ」プルダウンから「タスクの実行」を選択
3. タスクの設定として以下の選択を行う
   - 「環境」
     「既存のクラスター」から ECS タスク実行クラスターに選択
   - 「ネットワーキング」
     - 「VPC」: CDK で作成した VPC スタック(app-stg-**vpc**)に変更
     - 「サブネット」: public サブネット以外を削除
     - 「セキュリティグループ」: ECS タスクと同じセキュリテイグループ(app-{stg,prod}-**ecs-sg**)を設定
4. 最下部の「作成」ボタンを押下してタスクを実行する

### デプロイ

- 検証環境デプロイ

  ```sh
  npm run deploy-stg
  ```

- リソースの削除

  ```sh
  # 検証環境
  npm run destroy-stg
  ```

### phpMyAdmin を用いたレコードの管理(本番環境 or 検証環境)

1. AWS コンソールにて、踏み台 EC2 が起動していることを確認

2. 踏み台 EC2 に ssh 接続

   ```sh
   # 検証環境の場合。本番環境はAWSコンソールで確認
   ssh -i "e-privado-stg-bastion-key-pair.pem" ubuntu@ec2-57-180-30-253.ap-northeast-1.compute.amazonaws.com
   ```

   秘密鍵は パラメータストアに保存されている

3. phpMyAdmin の起動

   ```sh
   sudo docker run --name phpmyadmin -d --restart always -e PMA_HOST=DBのエンドポイント -p 8080:80 phpmyadmin/phpmyadmin
   ```

   ※ 注意
   * DB のエンドポイントは シークレットマネージャに保存されている

   * 前回起動時の phpMyAdmin コンテナが残っていてうまく起動できないことがあるので、その際は停止済みのコンテナを削除してから再試行

   * IpAddressToDBClientを設定していない場合、あらゆるユーザが閲覧できるようになるので、使用しないときはインスタンスを切る。また、本番環境は原則IPを指定すること。
