# BariPort

JAPAN CONNECT HACKTHON 用の、バックエンドリポジトリ

## 環境構築・実行方法

### 1. 必要な値の取得

AWS の操作を行うために、開発者から以下の２つの値の取得が必要です。

- AWS の IAM ユーザーの取得。
- AWS のアクセスキーの取得。

### 2. AWS CLI を開発用の PC にインストール

```
sudo pip install awscli  //pipを使用する場合
or
brew install awscli      //Homebrewを使用する場合
```

### 3. AWS CLI で、credentials を設定

```
aws configure --profile {任意のcredential名}
AWS Access Key ID [None]: ←取得したアクセスキーIDを設定
AWS Secret Access Key [None]: ←取得したシークレットアクセスキーを設定
Default region name [None]: ap-northeast-1 ←東京リージョンを設定
Default output format [None]: json ←出力形式をJSONに設定
```

### 4. リポジトリの Clone

### 5. モジュールのインストール

以下のコマンドで必要なモジュールをインストールします。

```
cd bari_port_back
npm i
```

## デプロイ方法

デプロイが必要な場合は、以下のコマンドを実行します。  
環境構築・実行方法の前準備が完了している必要があります。

```
AWS_PROFILE={credential名} pnpm sst deploy --stage=prod
```

## Swagger UI の起動方法

### 1. モジュールのインストール

以下のコマンドで必要なモジュールをインストールします。

```
cd swagger
npm i
```

### 2. Swagger UI の起動

以下のコマンドで Swagger UI 用のサーバを起動します。

```
npm start
```

### 3. ブラウザでアクセス

ブラウザで、`localhost:8080`にアクセスします。  
（2 の手順で自動で表示されることもあります）
