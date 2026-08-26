---
title: domain driven development
---

# 1 章 ドメイン駆動開発とは

### 用語説明

- ドメインとは対象とするシステム・サービスの**領域**のこと
  - 対象とするシステムによって大きく異なる
  - 例）
    - 動画配信システムなら、プレイヤー、動画、など
- モデルとは現実の概念や事象を抽象化した概念
  - その作業をモデリングという
- **ドメイン概念**をモデリングして得られたモデルを**ドメインモデル**
- **ドメインモデル**をコードに落とし込んだものを**ドメインオブジェクト**

### ドメイン駆動開発

- **業務上の知識(ドメイン概念)をソースコード(ドメインオブジェクト)に落とし込む開発手法**
- メリット
  - ドキュメントがなくてもソースコードを見れば仕様が把握可能
  - 可読性の向上
  - 変更に強くなる
- どうすれば？
  - 利用者を取り巻く世界を知ることが大事

---

# 2 章 値オブジェクト

- システム固有の値はプリミティブ型(`int`, `string` など)で表現することはできない
- 値をオブジェクトとして管理すること(**値オブジェクト**)でシステム固有の値として保持する

```go
// FullNameという値をオブジェクトとして扱う
type FullName struct {
	lastName  string
	firstName string
}

func NewFullName(lastName, firstName string) FullName {
	return FullName{
		lastName:  lastName,
		firstName: firstName,
	}
}
```

## 値の性質

1. 不変である
2. 交換が可能である
3. 等価性によって比較される

### 1. 不変である

- 途中で値の変更を許すと予期せぬ時に書き変わってバグを生む
- 値自身が自身を変更する振る舞いをもつのはおかしい
- 値を変更したい場合、新たにインスタンスを生成する

```go
// bad
fullName := NewFullName("yamada", "taro")
fullName.ChangeLastName("sato")
```

### 2. 交換が可能である

- 値はそれ自身は変更できないが交換は可能である
- しかし、fullName の値が再代入されるため、バグの元になるためあまりやらない方が良い気がする（？）

```go
fullName := NewFullName("yamada", "taro")
fullName = NewFullName("suzuki", "jiro")
```

### 3. 等価性によって比較される

- 値オブジェクトは値なので、以下のようにフィールドを取り出して比較するコードは不自然
- 下のようにすると例えば、ミドルネームというフィールドが追加されても修正は `Equal` メソッドだけで済む

```go
// bad
userA := NewFullName("yamada", "taro")
userB := NewFullName("yamada", "taro")
if userA.lastName == userB.lastName && userA.firstName == userB.lastName {
	fmt.Println("同一")
}
```

```go
// good
func (f FullName) Equal(other FullName) bool {
	return f.lastName == other.lastName && f.firstName == other.firstName
}

userA := NewFullName("yamada", "taro")
userB := NewFullName("yamada", "taro")
userA.Equal(userB)
```

## 業務上のルールを値オブジェクトに含める

- ルール
  - 名前は 1 文字以上

```go
// FullNameという値をオブジェクトとして扱う
func NewFullName(lastName, firstName string) (FullName, error) {
	if lastName == "" {
		return FullName{}, errors.New("苗字は1文字以上で指定する必要があります")
	}
	if firstName == "" {
		return FullName{}, errors.New("名前は1文字以上で指定する必要があります")
	}
	return FullName{lastName: lastName, firstName: firstName}, nil
}
```

## どこまで値オブジェクトにすべきか

- 以下２点を基準に値オブジェクトにする
  - ルールが存在するか
  - それ単体で取り扱いたいか

## 値オブジェクトに振る舞いを持たせる

- オブジェクトに対する処理を振る舞いとして一箇所にまとめることで、自身に関するルールを語るドメインオブジェクトらしさを帯びる
  - 変更が容易
  - ドキュメント不要
- 値に関する処理は全て値オブジェクト内に書く

```go
type Money struct {
	amount int64
}

func (m Money) Add(other Money) Money {
	return Money{amount: m.amount + other.amount} // 値は不変のため新たに値を返す
}
```

## 値オブジェクトを利用するメリット

1. 表現が増す
2. 不正な値を存在させない
3. 誤った代入を防ぐ
4. ロジックの散在を防ぐ

### 1. 表現が増す

- string 型や int 型では文字列、数値であれば何でも許容するため、文字列、数値という意味しかもたない
- 業務上で扱う値は、必ず何かしらのルールや、意味を持つ
- 値オブジェクトにルールを記載するので、仕様の代わりになり、何を表しているのかが明確になる

```go
// bad
productCode := "a-132-2012"

// good
type ProductCode struct {
	area        string
	auth        string
	createdYear string
}
```

### 2. 不正な値を存在させない

- システムで扱う値は必ずルールが存在する
  - 固定の電話番号は 10 桁
  - ユーザーは 3 文字以上 10 文字以下
  - メールアドレスはアルファベットのみ、任意の文字列@ドメイン名
- プリミティブの string 型だと文字列であれば何でも許容される

```go
// 11であるが許容される
phoneNumber := "09232322344"
```

- 値を利用する色んな箇所でチェックする必要がある
- 仕様変更になった時に全てを修正する必要があり大変

```go
if len(phoneNumber) == 10 {
	// 正しい挙動の処理
} else {
	return errors.New("電話番号が不正です")
}
```

### 3. 誤った代入を防ぐ

- プリミティブ型で扱うとエラーが出ずにバグの原因となる
- id と name はそれぞれ値オブジェクトにすると良い

```go
// エラーが出ない
func createUser(name string) User {
	id := name
	return NewUser(name, id)
}
```

### 4. ロジックの散在を防ぐ

- インスタンスプロパティをそのまま返す getter を用意することは外部に同じようなロジックが散在する可能性がある
- setter も同様
- setter, getter を利用したい場合は参照先のロジックが値オブジェクトに移動できないか検討
- なるべく getter, setter は利用しない

## まとめ

- 値オブジェクトはシステム固有の値を作ること
- システムで取り扱う値には必ずルールがあるため、プリミティブ型で取り扱わないようにする

---

# 3 章 ライフサイクルのある「エンティティ」

- ドメイン駆動におけるエンティティは、**値オブジェクトと同じくドメインモデルを実装したドメインオブジェクト**である
- 値オブジェクトとの違いは、**同一性**によって識別されるか否か
- 例
  - 値オブジェクト(氏名オブジェクト)
    - 異なる値にすれば別のもの
    - 同じ値にすれば同じもの
  - エンティティ(ユーザオブジェクト)
    - 氏名が途中で変わっても同一
    - 同じ氏名でも異なる人物である可能性がある
    - **固有の同一性(identity)を保つ**

## エンティティの性質

1. 可変である
2. 同じ属性であっても区別される
3. 同一性(id)によって区別される

### 1. 可変である

- エンティティは可変なオブジェクトである
- ユーザ名(ニックネーム)は頻繁に変えたりできる
- 制限をかけて可変にする
  - そのままの値をセットするセッターはダメ

```go
type User struct {
	name Name
}

func (u *User) ChangeName(name Name) error {
	if name == (Name{}) {
		return errors.New("名前が空です")
	}
	u.name = name
	return nil
}
```

### 2. 同じ属性であっても区別される

- 値オブジェクトは、値が同じものであれば等価性の原理によって全く同じものと表せる
- エンティティはユーザーのように同じ氏名でも同じではない
- **ID**(identity)によって区別される

```go
type User struct {
	id   UserID // ユーザを識別する固有のID
	name Name
}

func NewUser(id UserID, name Name) (*User, error) {
	if id == (UserID{}) {
		return nil, errors.New("idが必要です")
	}
	if name == (Name{}) {
		return nil, errors.New("nameが必要です")
	}
	return &User{id: id, name: name}, nil
}
```

### 3. 同一性によって区別

- 固有識別 id により識別されるため、再代入すべきではない
- 同一性を比較する処理が必要

```go
func (u User) Equal(other User) bool {
	// 同一性は id で判定する
	return u.id == other.id
}
```

## エンティティにする判断基準

- ライフサイクルがあるかどうか
- 同一性で判別したいかどうかも大事
- 例
  - ユーザの場合は、アカウント作成により生を受けて、削除されて死を迎える
  - つまり、ライフサイクルがあるのでエンティティ
- 迷ったら一旦**値オブジェクト**にする
  - 不変性を保つことでコードをシンプルにする

## ドメインオブジェクトにするメリット

1. コードのドキュメント性が高まる(可読性)
2. ドメインの変更がコードに伝えやすくなる(保守性)

### 1. コードのドキュメント性が高まる

- プロジェクト途中参画者がプロジェクトのことを知るにはドキュメントが必要
  - ドキュメントは更新されていない可能性あり
  - 抜け漏れがある可能性あり
- ドメインオブジェクトにすればコードを見るだけで、仕様を把握することができる

```go
// ドメインオブジェクトではない場合
type User struct {
	name string // 名前はstringの情報しかわからん
}

func NewUser(name string) User {
	return User{name: name}
}

func (u User) Name() string {
	return u.name
}

// ドメインオブジェクト
type User struct {
	name string
}

func NewUser(name string) (User, error) {
	if len(name) < 3 {
		return User{}, errors.New("名前は3文字以上で指定してください")
	}
	return User{name: name}, nil
}

func (u User) Name() string {
	return u.name
}
```

### 2. ドメインの変更がコードに伝えやすくなる

- ドメインの変更(仕様変更)が、コードに反映させやすい
- 例
  - ユーザ名の仕様が 3 文字以上から 5 文字以上になってもユーザオブジェクトを参照するだけで OK

## まとめ

- 豊かな振る舞いを持たせることで、曖昧さを省く

---

# 4 章 不自然さを改善するドメインサービス

- ドメインをそのままコードに落とし込んだ時に、値オブジェクトやエンティティの振る舞いとして定義するのは不自然と感じる場合がある
- それを解決するために、「ドメインサービス」がある
- ドメイン(値オブジェクト、エンティティ)のためのサービス
- **ドメインオブジェクトに書くことが不自然な場合のみドメインサービスを利用**

## 不自然な振る舞いとは

- ユーザの重複を確認することを考える
- エンティティに記述するとユーザー自身に自身の重複を確認するという現実世界ではおかしな振る舞いとなる

```go
user := NewUser(userID, userName)
// user自身に重複を確認するのはおかしい
duplicateCheckUser := user.Exists(user)
```

## 不自然さを解決するオブジェクト

- 自身の振る舞いを示さないドメインのためのドメインサービスを定義することで解決

```go
type UserService struct{}

func (s UserService) Exists(user User) bool {
	// 重複を確認する処理
	return false
}

userService := UserService{}
user := NewUser(userID, userName)
// 不自然さはなくなった
duplicateCheckUser := userService.Exists(user)
```

## ドメインサービスの濫用がもたらすもの

- ドメインオブジェクトが情報を持たないものとなり、ルールや仕様がわからなくなる
- ドメイン駆動の本来の志向とは真逆のものとなってしまう
- ドメインオブジェクトかドメインサービスどっちに書くか迷った場合はドメインオブジェクトに書くこと
- **なるべくドメインサービスは利用しないようにする**

## ドメインサービスの命名

- プログラマーにドメインサービスと分かるようにする必要がある
  - ドメイン概念+Service

## まとめ

- 値オブジェクトやエンティティに記述するとドメイン的に不自然な振る舞いになる場合は、ドメインサービスを利用する
- ただし不自然となる場合に限定し、なるべく利用しないように心がける

# 5 章 データにまつわる処理を分離する「リポジトリ」

## リポジトリとは

- データストア(データベース)とドメインオブジェクトの橋渡し
- データの**永続化**と**再構築**だけを担う

ドメインオブジェクト ⇄ リポジトリ ⇄ データストア

- ドメインオブジェクトの値をリポジトリに渡すのは、アプリケーションサービスの役割

## データストアの責務

- リレーショナルデータベースを処理するコードは、ややこしい記述になりがち
  - ぱっと見てどのような処理を施しているのかがわからない
- ドメインにとって、何を使ってどのようにデータを保存・復元するのかは重要ではない
- ドメインにとって大事なのは、**あるデータを保存する(永続化)、また、あるデータを取得する(再構築)すること**

```go
type UserApplicationService struct {
	// ユーザリポジトリのインターフェースに依存する
	userRepository UserRepository
}

func NewUserApplicationService(userRepository UserRepository) *UserApplicationService {
	return &UserApplicationService{userRepository: userRepository}
}

// ユーザを作成して保存する（永続化）
func (s *UserApplicationService) CreateUser(userName string) error {
	// ユーザ作成
	user, err := NewUser(NewUserName(userName))
	if err != nil {
		return err
	}
	// 保存（永続化）
	return s.userRepository.Save(user)
}
```

```go
type UserService struct {
	// ユーザリポジトリのインターフェースに依存する
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

// ユーザが存在するか判断する（再構築した結果で判定する）
func (s *UserService) Exists(user User) bool {
	// ユーザ名からユーザを再構築する
	found, err := s.userRepository.FindByName(user.name)
	if err != nil {
		return false
	}
	return found != nil
}
```

- 大事なのはリポジトリは、データの**永続化と再構築だけ**を担う
- 以下のように、ユーザを引数として受け取り、存在するかどうかを返す Exists メソッドをリポジトリにはおかない
- ユーザが存在するかどうかは、ドメイン側の知識なのでドメインサービスに持たせるべき

```go
type UserRepository interface {
	// 永続化なので◯
	Save(user User) error
	// ユーザ名よりユーザ情報を返すデータの再構築なので◯
	FindByName(userName UserName) (*User, error)
	// 存在判定だけで再構築していないため❌
	Exists(userName UserName) bool
}
```

## リポジトリクラス

- 各データベースに保存するための SQL などを記載
- データベース特有の記述をして OK

```go
type MySQLUserRepository struct{}

func (r *MySQLUserRepository) FindByName(name UserName) (*User, error) {
	// database特有の書き方
	return nil, nil
}

func (r *MySQLUserRepository) Save(user User) error {
	// database特有の書き方
	return nil
}
```

## 利用

```go
userRepository := &MySQLUserRepository{}
svc := NewUserApplicationService(userRepository)
if err := svc.CreateUser("taro"); err != nil {
	log.Fatal(err)
}
```

## リポジトリのインターフェース

- リポジトリはインターフェースに依存させるべき
- 途中でデータベースの種類を変更しても、インターフェースにさえ依存させておけば、利用側のコード(ドメイン側)は変更しなくて済む
- テストの際もテスト用のリポジトリに変更するだけで済む

```go
type UserRepository interface {
	Save(user User) error
	FindByName(userName UserName) (*User, error)
}
```

## テスト

- テスト用にデータベースを用意して、テストごとにデータテーブルを用意するのは手間と時間がかかる
- 実際にデータベースを用意するのではなく、メモリ上をデータベースに見立ててテストすると気軽にテストを実施できる

```go
type InMemoryUserRepository struct {
	store map[string]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		store: map[string]User{
			"4": {name: "taro", id: "4"},
			"2": {name: "jiro", id: "2"},
		},
	}
}

func (r *InMemoryUserRepository) FindByName(userName UserName) (*User, error) {
	for _, user := range r.store {
		if userName == user.name {
			cloned := r.clone(user)
			return &cloned, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) Save(user User) error {
	r.store[user.id.String()] = r.clone(user)
	return nil
}

// 値型なので代入でコピーになる
func (r *InMemoryUserRepository) clone(user User) User {
	return user
}
```

- テスト時

```go
func TestCreateUser(t *testing.T) {
	repo := NewInMemoryUserRepository()
	svc := NewUserApplicationService(repo)
	if err := svc.CreateUser("taro"); err != nil {
		t.Fatal(err)
	}
	got := repo.store["1"]
	if got.name != "taro" {
		t.Errorf("unexpected user: %+v", got)
	}
}
```

## オブジェクトリレーショナルマップでリポジトリ構築

- ORM とは
  - オブジェクト指向パラダイムを使用してデータベースからデータを取得及び操作するのに役立つ手法
- メリット
  - MVC になるためコードがクリーン
  - SQL クエリを作成する必要がない
- デメリット
  - 複雑なクエリでのパフォーマンスの問題
- Go

  - GORM
  - sqlx
  - ent
  - sqlc

- まずは、SQL 文をしっかり知ることが大事(?)

## リポジトリの振る舞い

### 永続化を持つ

- 値を更新するものはリポジトリではない

```go
// ダメ
type UserRepository interface {
	UpdateName(id UserID, name UserName) error
}
```

- 値の更新はオブジェクトに持たせる

```go
func (u *User) ChangeName(name UserName) error {
	u.name = name
	return nil
}
```

- 削除はリポジトリの責務

```go
// ok
type UserRepository interface {
	Delete(user User) error
}
```

### 再構築

- 検索条件ごとにメソッドを分ける

```go
type UserRepository interface {
	FindByID(id UserID) (*User, error)     // ID で再構築する
	FindByName(name UserName) (*User, error) // 名前で再構築する
}
```

## まとめ

- データベースを触る部分はややこしくなるが、ドメインにとって大事ではないため、リポジトリにまとめることで可読性をあげる
- リポジトリはインターフェースに依存させることで途中で DB が変わってもリポジトリを変えるだけで、アプリケーションサービス、値オブジェクトなどは変更しなくても良い

# 6 章 アプリケーションサービス

- **リポジトリとドメインオブジェクトを利用して、ユースケースを実現するためのオブジェクト**

## ユースケース

- ユーザ機能を例に以下のユースケースを考える
  - ユーザ登録
  - ユーザ情報確認
  - ユーザ情報を更新
  - 退会
- ユーザ情報
  - エンティティ
  - 保持する値
    - ユーザ名
      - 値オブジェクト
    - ユーザ ID
      - 値オブジェクト
- 重複を確認する処理
  - ドメインサービス
- リポジトリ
  - ユーザ情報を検索
  - ユーザ情報を作成
  - ユーザ情報を削除

![ユーザ機能のクラス図](./diagrams/ch06-usecase.svg)

## 1. ユーザを登録する

- ユーザの重複がないかを確認し、重複がなければ登録

![ユーザ登録](./diagrams/ch06-register-user.svg)

## 2. ユーザ情報を取得

- ユースケース詳細
  - ユーザ ID よりユーザを検索し、存在したらユーザ情報を返す
- User エンティティをそのまま返却してはならない
  - **ドメインオブジェクト(ドメインサービス)の振る舞いを呼び出すのは、アプリケーションサービス内だけにとどめるのが大事**
    - 依存が大きくなり、ドメインオブジェクトの変更が難しくなる
    - アプリケーションサービス以外にロジックが広がる可能性がある
- ユーザ情報は DTO を利用して返す
  - Data Transfer Object はデータを転送するためだけのオブジェクト
  - User 情報を DTO に移し替えて返すことで、ドメインオブジェクトはアプリケーションサービス以外からアクセスさせない

![ユーザ情報の取得](./diagrams/ch06-get-user.svg)

- UserData のコンストラクタに渡す引数はドメインオブジェクトのインスタンスそのものを渡す
  - ドメインオブジェクトのインスタンスプロパティが増えても、修正が DTO(UserData)のみで済む
  - DTO は値を転送するためだけのオブジェクトなので、ドメインオブジェクトの振る舞いを呼ばないようにする

## 3. ユーザ情報を更新する

- ユースケース

  - ユーザ情報を一括で更新する
    - 今後メールアドレス、住所、電話番号などの情報が加わることを想定して作成する

- アップデートする各データを引数に渡すのは良くない
  - アップデートする情報が増えるたびに、シグネチャ(引数)が増えるのは良くない
    - 冗長になる
    - 利用側も変更しないといけないのがダメ
    - 実装側の中身だけを変えるだけで済む
- リポジトリの save メソッドは SQL を使うのであれば `upsert` を利用

![ユーザ情報の更新](./diagrams/ch06-update-user.svg)

## 4. 退会処理

- ユースケース
  - アカウントが不要になったら退会を実施する

![退会処理](./diagrams/ch06-withdraw.svg)

## アプリケーションサービスにドメインオブジェクトのルールを持ってこない

- ドメインルールをアプリケーションサービスに持ってくると、同じロジックがアプリケーションサービス内に散在し、仕様変更の際に変更箇所が多くなり、抜け漏れが生じてバグとなる

## 凝縮度の高いアプリケーションサービス

- 凝縮度とは、モジュールの責任範囲の集中度合い
- 凝縮度の高いモジュールは、再利用性、可読性、堅牢性、信頼性が高いため、なるべく凝縮度を上げるべき
- インスタンス変数は全てのメソッドで利用するモジュールは凝縮度が高い

![ユースケースごとのアプリケーションサービス](./diagrams/ch06-cohesion.svg)

- ユーザに関するモジュール処理はパッケージとしてまとめる
  - パッケージはディレクトリ単位に切り分けて実現
    - application/user/

![アプリケーションサービスのパッケージ構成](./diagrams/ch06-package.svg)

## アプリケーションサービスのインターフェース

- アプリケーションサービスをインターフェースにすることで、クライアント側とアプリケーション側は互いに分担して作業ができる
- クライアント側は、Web 開発で考えたらフロントエンドのことを指す

## サービスとは？

- クライアントのために、何かを行うものであり、自身の振る舞いを持たない
  - ドメインサービス：ドメインのための振る舞いを記述
    - ドメインは、対象とするシステム・サービスの領域
  - アプリケーションサービス：アプリケーションのための振る舞いを記述
    - アプリケーションは、ユースケース・機能
- **サービスは状態を持たない**
  - 自身の振る舞いを変更するインスタンス変数は持たない
  - 自身の振る舞いを変化させないインスタンス変数は OK

## まとめ

- アプリケーションサービスでは、ユーザのユースケースをドメインオブジェクトとリポジトリより、ボトムアップ形式で組み立てる
- ドメインのルールがアプリケーションに入っていないかをよく確認することが大事

# 7 章 依存関係のコントロール

## 依存関係とは

- あるオブジェクトに対して、他のオブジェクトが存在しないと成り立たない関係性のこと

```go
// ObjectAはObjectBが存在しないと成り立たない
// ObjectAはObjectBに依存している
type ObjectA struct {
	objectB ObjectB
}
```

![依存関係](./diagrams/ch07-dependency.svg)

## インターフェースに依存

- インターフェースに対しても依存が成り立つ
  - インターフェースを実装することを**実現**という

![インターフェースの実現](./diagrams/ch07-interface.svg)

- 利用側はインターフェースに依存させ、実装もインターフェースをもとに実現する
  - 利用側は特定のオブジェクトに依存しなくなる
  - 利用者側(アプリケーションサービス側)はコードを変更せずにリポジトリを変更できるようになる
  - リポジトリ実現する人と、アプリケーションサービスを作成する人はお互いの実装を気にせず実装できる

![複数のリポジトリ実装](./diagrams/ch07-multiple-impl.svg)

## 依存関係逆転の原則

- 原則

  - 上位レベルのモジュールが下位レベルのモジュールに依存してはならない
    - どちらのモジュールも抽象に依存すべき
  - 抽象は詳細(下位レベル)に依存してはならない
    - 抽象の主導権は上位レベルに持たせる

- 上位レベルはよりドメインに近いもの
  - アプリケーションサービスとデータストアを扱うリポジトリの実装を比べるとアプリケーションサービスの方が上位レベル

![依存関係逆転の原則への違反](./diagrams/ch07-dip-violation.svg)

![依存関係逆転の原則](./diagrams/ch07-dip.svg)

- インターフェースの主導権は上位レベルに持たせ、下位レベルはそれに従う

![インターフェースの主導権](./diagrams/ch07-interface-ownership.svg)

## IoC Container パターン

- リポジトリを利用するアプリケーションサービスは、リポジトリをコンストラクタの引数で受け取る

```go
type UserApplicationService struct {
	// ユーザリポジトリのインターフェースに依存する
	userRepository UserRepository
}

func NewUserApplicationService(userRepository UserRepository) *UserApplicationService {
	return &UserApplicationService{userRepository: userRepository}
}

// 本番
svc := NewUserApplicationService(NewMySQLUserRepository())

// テスト
svc := NewUserApplicationService(NewInMemoryUserRepository())
```

- コンストラクタ内で実装を new すると、差し替えのたびにアプリケーションサービスを直すことになる

```go
func NewUserApplicationService() *UserApplicationService {
	return &UserApplicationService{
		userRepository: NewMySQLUserRepository(), // 実装にべったり依存
	}
}
```

- Service Locator は依存が分かりにくく、テストが壊れやすい

```go
func NewUserApplicationService() *UserApplicationService {
	userRepository, _ := ServiceLocator.Get("UserRepository")
	return &UserApplicationService{userRepository: userRepository.(UserRepository)}
}
```

- 規模が大きくなったら DI コンテナも選択肢になる

```go
func NewUserApplicationService(userRepository UserRepository) *UserApplicationService {
	return &UserApplicationService{userRepository: userRepository}
}

container := dig.New()

// 依存関係の登録
container.Provide(NewMySQLUserRepository)
container.Provide(NewUserApplicationService)

// インスタンス化
var svc *UserApplicationService
if err := container.Invoke(func(s *UserApplicationService) {
	svc = s
}); err != nil {
	log.Fatal(err)
}
```

## まとめ

- 依存関係を何も考えずに設計すると、修正しにくく、拡張性がないソフトウェアになる
- 依存関係をコントロールし、柔軟に保つことが大事

---

# 8 章 ソフトウェアシステムを組み立てる

- アプリを利用するにはユーザが実際に触れる UI が必要
- UI には CLI と GUI の２種類ある
- CLI
  - コマンドラインインターフェース
  - コマンドプロンプトに対して文字列の命令を処理をする
  - `$ cd .`
- GUI
  - グラフィックユーザインターフェース
  - コマンドの知識がなくても直感的に処理をする

## CLI に組み込む

![CLI への組み込み](./diagrams/ch08-cli.svg)

![CLI のシーケンス](./diagrams/ch08-cli-sequence.svg)

## MVC フレームワークに組み込む

![サーバー起動時の依存関係登録](./diagrams/ch08-startup.svg)

- コントローラー(Controller)とクライアント(View)とアプリケーションサービス(Model)の関係

![コントローラーとアプリケーションサービス](./diagrams/ch08-mvc.svg)

![ユーザ登録リクエストのシーケンス](./diagrams/ch08-mvc-sequence.svg)

## ユニットテスト

- テスト観点
  - 境界値テスト
    - 正常
      - データが保存されているか
    - エラー
      - エラーが正しく表記されるか

## まとめ

- ユーザインターフェースによらず、アプリケーションサービスは不変
  - アプリケーションサービスのテストが可能になる
  - ユーザインターフェースは変更しやすい

# 9 章 複雑な生成処理を行うファクトリ

- オブジェクトの生成は複雑になりがち
- その責務を行うファクトリオブジェクトを考える

## ユーザ ID の採番処理

- 採番テーブルの利用を考える
  - ユーザドメインにデータベースの採番を持ってくると、ドメインが見えにくくなり、可読性が下がる

```go
type User struct {
	userID UserID
	name   UserName
}

func NewUser(name UserName) (*User, error) {
	if name == (UserName{}) {
		return nil, errors.New("名前が必要です")
	}
	// データベースの処理
	db, err := sql.Open("mysql", "user:password@/")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 何かSQLの操作
	var sequenceID string
	err = db.QueryRow("SELECT LAST_INSERT_ID(id + 1) FROM seq").Scan(&sequenceID)
	if err != nil {
		return nil, err
	}

	return &User{
		userID: NewUserID(sequenceID),
		name:   name,
	}, nil
}
```

- ユーザを作成するためのファクトリを作成する

```go
// インターフェースを作成することで異なるDB利用となってもすぐ変更可能
type UserFactory interface {
	Create(name UserName) (*User, error)
}

// ユーザ作成をするためのファクトリ
type MySQLUserFactory struct{}

func (f *MySQLUserFactory) Create(name UserName) (*User, error) {
	db, err := sql.Open("mysql", "user:password@/")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 何かSQLの操作
	var sequenceID string
	err = db.QueryRow("UPDATE seq SET id = LAST_INSERT_ID(id + 1)").Scan(&sequenceID)
	if err != nil {
		return nil, err
	}
	// uuidを利用する場合は、sequenceID := uuid.NewString() でもOK
	userID := NewUserID(sequenceID)
	return NewUser(name, userID)
}

// ユーザオブジェクトはシンプルになる
type User struct {
	userID UserID
	name   UserName
}

func NewUser(name UserName, userID UserID) (*User, error) {
	if name == (UserName{}) {
		return nil, errors.New("名前が必要です")
	}
	if userID == (UserID{}) {
		return nil, errors.New("idが必要です")
	}
	return &User{userID: userID, name: name}, nil
}
```

```go
// 利用
userFactory := &MySQLUserFactory{}
// ファクトリー経由でユーザ作成
user, err := userFactory.Create(userName)
```

- 実際のデータベースを触らない場合は、UserFactory を満たすテスト用の実装を作るだけで OK

- ファクトリの存在を気づかせるディレクトリ構成

```
./domain/user/user.go
./domain/user/factory.go
./infrastructure/mysql/user_factory.go
./infrastructure/memory/user_factory.go
```

## DB の自動採番の利用

- 自動採番の利用には注意が必要

  - ユーザを識別する値なので保存するまで採番されないのは不自然
  - 保存してから userID をユーザ型にセットするため、セッターが必要
    - セッターは不本意に値が変えられる可能性があるため利用すべきでない

  **ルールを決める必要がある**

## ファクトリメソッド

- あるユーザ(オーナー)がグループを作ることを考える

```go
type User struct {
	userID UserID
}

func (u User) ID() UserID {
	return u.userID
}

// 利用側
user, _ := userFactory.Create(name)
group := NewGroup(
	user.ID(), // userIDをゲッターより取得
	NewGroupName("サーフィン"),
)
```

- ゲッターを利用している箇所を元の型へ移せないか検討すると、ゲッターが不要になる
  - 依存度が小さく
  - より高いレベルのルールがわかる
  - **セッター、ゲッターを利用している処理は、利用元の型に移動できないか検討**

```go
type User struct {
	userID UserID
}

func (u User) CreateGroup(name GroupName) Group {
	return NewGroup(u.userID, name)
}

// 利用側
group := user.CreateGroup(NewGroupName("サーフィン")) // あるユーザがグループを作成するというドメインモデルの高いレベルの仕様が伺える
```

## ファクトリを利用するかの検討

- コンストラクタが複雑な処理の場合はファクトリなどに任せられないかを検討する
- 特に、コンストラクタでオブジェクトを生成している箇所は、ファクトリを使えないかを検討する
  - そのオブジェクトが変更されるたびに修正しないとダメ

## まとめ

- 複雑な処理をファクトリなどによりカプセル化することでより明瞭、柔軟なコードになる
- セッター、ゲッターを利用している処理は、利用元の型に移動できないか検討
  - セッター、ゲッターもなるべく利用しない
- コンストラクタが複雑な場合、ファクトリがないかを検討

# 10 章 データの整合性を保つ

## 重複データを保存しないルールへの対応

- ユーザ名の重複をさせないルールを考える
- 同時タイミングで同一名前のユーザ登録をするとユーザ名が重複する可能性がある
- 解決策として DB の**ユニークキー制約**を利用する
  - 同じ名前のデータは保存できなくなる
- これを使うと、ユーザ名の重複を確認するコード(ドメインサービス)に書かなくてもよくなる
  - **コードを見ただけでユーザ名が重複してはいけないというルールが見えなくなる**
  - コードをドキュメントにするというドメイン駆動設計の思想から外れる
- **ユニーク制約とユーザ名の重複を確認するコード(ドメインサービス)の両方を書くことが大事**

## 処理が途中で失敗した場合のデータ不整合

- EC サイトで商品を購入する際の以下の例を考える

  1. 在庫を DB に確認し、OK だった
  2. ポイント残高減らし、DB 反映
  3. 在庫を減らし、DB 反映
  4. 注文処理完了

- この際、2 の間に在庫がなくなった場合、ポイントは減るが、商品が購入できない事態が発生する可能性がある
- それを阻止するために**トランザクション**を利用する

  1. 在庫を DB に確認し、OK だった
  2. ポイント残高減らし、一時保存
  3. 在庫を減らし、一時保存
  4. **コミット(DB にポイント残高、在庫を反映)**
  5. 注文処理完了

- 方法
  - スコープを利用したパターン
    - C#
  - AOP を利用したパターン
    - TypeScript, Java, etc
  - 関数のラップを利用したパターン
    - Go
  - ユニットオブワークを利用したパターン

```go
// 関数のラップを利用した方法

// ターゲット関数
func myFunction() {
	fmt.Println("myFunction")
}

// アドバイス
func myAdvice() {
	fmt.Println("myAdvice")
}

// 関数をラップして after 相当の処理を付ける
func withAfter(fn func(), after func()) func() {
	return func() {
		fn()
		after()
	}
}

proxyFunction := withAfter(myFunction, myAdvice)

// 関数呼び出し
proxyFunction()
```

# 11 章 アプリケーションを 1 から組み立てる

![サークルのユースケース](./diagrams/ch11-usecase.svg)

## サークルを作成する

![サークル作成の構成](./diagrams/ch11-circle-create.svg)

# 12 章 集約

- 全体(集約ルートのオブジェクト)-部分(集約境界内のオブジェクト)関係を表す
  - 全体に当たる型が部分の型を包含する
- 外部から集約の境界内のオブジェクトの操作は、集約ルートを通して行う
- 集約のうち、集約境界内オブジェクトの操作はその集約ルートを通してのみ行う場合**コンポジション**と言い、菱形の黒塗りで表す

![User 集約](./diagrams/ch12-aggregate-user.svg)

- 集約ルートオブジェクトの集約は、Circle オブジェクトのみで User オブジェクトを操作しないので、ただの集約と表記する

![Circle 集約](./diagrams/ch12-aggregate-circle.svg)

## デメテルの法則

オブジェクト同士のメソッドを呼び出す際の秩序を守るためのガイドライン

- メソッドを呼び出すオブジェクトは以下の 4 つ

  1. オブジェクト自身
  2. 引数として渡されたオブジェクト
  3. インスタンス変数
  4. 直接インスタンス化したオブジェクト

- 例

```go
// OK
type Circle struct {
	members []User
}

func (c *Circle) Join(user User) {
	// 3. インスタンス変数のメソッド呼び出し
	c.members = append(c.members, user)
}

circle := &Circle{}
// 1. オブジェクト自身
circle.Join(user)

// デメテルの法則に違反するコード例 NG
// オブジェクトのインスタンス変数のメソッド
_ = len(circle.members)
circle.members = append(circle.members, user)
```

## ドメインのルールは全てドメインオブジェクトにする

- ゲッターをなるべく利用しない
- ゲッターで値だけを取得している箇所では、その値を利用してロジック(ルール)を記載する恐れあり
  - そのロジック(ルール)を大元の型(ドメインオブジェクト)にもってきて、ゲッターを利用しない
  - そうすることでデメテルも自然と守られる傾向がある
- **ルールはドメインオブジェクトに収める**

## 通知オブジェクトで情報を隠蔽する

- ゲッターがないと困ることがある
- ドメインオブジェクトのデータを永続化するためにはゲッターを利用している

```go
type User struct {
	id   UserID
	name UserName
}

func (u User) ID() UserID {
	return u.id
}

func (u User) Name() UserName {
	return u.name
}

type MySQLUserRepository struct{}

func (r *MySQLUserRepository) Save(user User) error {
	// ゲッターで内部データを取り出して DB に保存する
	id := user.ID()
	name := user.Name()
	saveData := map[string]any{"id": id, "name": name}
	_ = saveData
	return nil
}
```

- それを解決するために通知オブジェクトを利用

```go
type UserNotifier interface {
	NotifyID(id UserID)
	NotifyName(name UserName)
}

type UserDataModelBuilder struct {
	id   UserID
	name UserName
}

func (b *UserDataModelBuilder) NotifyID(id UserID) {
	b.id = id
}

func (b *UserDataModelBuilder) NotifyName(name UserName) {
	b.name = name
}

func (b *UserDataModelBuilder) Build() map[string]string {
	// リポジトリが使うデータ形式に変換する
	return map[string]string{
		"id":   b.id.String(),
		"name": b.name.String(),
	}
}

// ユーザ型にゲッターは不要になる
type User struct {
	id   UserID
	name UserName
}

func (u User) Notify(n UserNotifier) {
	n.NotifyID(u.id)
	n.NotifyName(u.name)
}

type MySQLUserRepository struct{}

func (r *MySQLUserRepository) Save(user User) error {
	// 通知オブジェクトを渡して内部データを受け取る
	builder := &UserDataModelBuilder{}
	user.Notify(builder)
	saveData := builder.Build()
	_ = saveData
	return nil
}
```

- ただ、値オブジェクトを文字列にする処理は仕方ない

```go
type UserID struct {
	id string
}

func (id UserID) String() string {
	return id.id
}
```

## どの単位でリポジトリを作成するのか

- **集約の単位でリポジトリは作成する**
- 集約に対する変更、永続化の依頼も集約ごとに対して行うことが大事
  - そうしないと同じコードが散在する

## 集約の値で集約を保持する場合は制約をかける

- id のみ集約に含める
  - id さえあれば一意のデータを取得できる
  - メモリの節約
  - 不慮のメソッド呼び出しがなくなる

![Circle 集約に UserID だけ含める](./diagrams/ch12-aggregate-userid.svg)

## ID のゲッターの是非

- ID を利用してビジネスロジックを書くことはほとんどないため、ゲッターにしても問題になりにくい

## 集約の単位

- 集約の単位はなるべく小さく保つ、同一トランザクション内で多くの集約を扱わないようにする
  - 処理が多いと、トランザクションが失敗する可能性が高くなる
- トランザクションが失敗するとデータがロックされる可能性がある

## 言葉の齟齬はドメインオブジェクトで解決

- コードとドメインルールに齟齬がある場合は、ドメインオブジェクトに解決するためのメソッドを用意してあげて解決する

# 13 章 仕様(specification)

- ドメインオブジェクトの１種
- 評価を行うオブジェクト
- ドメインオブジェクト自身のメソッドで表すのに複雑な場合は仕様として抽出する
  - 複雑なとは、リポジトリが絡んだりする評価
- リポジトリが絡んだルールをドメインオブジェクト自身のメソッドに記述すると、ドメインという高レベルに低レイヤーのリポジトリを置くのは、良くないため、仕様として抽出する

## 仕様に対してもリポジトリの利用は控える

- 仕様はドメインオブジェクトの１種なので、リポジトリの利用は良くないという考え
- ファーストクラスコレクションパターンで解決
  - 特化した集合オブジェクト
- **ファーストクラスコレクションを利用すれば、仕様としてオブジェクトを定義しなくても、ドメインオブジェクトのメソッドとして、評価を行うメソッドを定義できる**
  - リポジトリを使わないため

```go
func (c Circle) IsFull() bool {
	premium := c.members.CountPremium(false)
	limit := 50
	if premium < 10 {
		limit = 30
	}
	return c.members.Len() >= limit
}
```

## 仕様とリポジトリの組み合わせ

- リポジトリでは、データの復元つまり検索をするが、その検索のルールがリポジトリに隠れる可能性がある
  - 例、おすすめのサークルを検索
    - 直近 1 カ月以内に登録
    - メンバーが 10 名以上
- このルールをリポジトリに記述するのではなく、仕様として定義することが大事
- **あくまでルールはドメインオブジェクトに記述することが大事**

### リポジトリに仕様を渡す方法

- 仕様を別途定義して、それをリポジトリが利用する形にするとルールをリポジトリに書かなくて済む

![仕様とリポジトリ](./diagrams/ch13-specification.svg)

- **この方法だと全サークルを取得してきて、合致するサークルを探すため、パフォーマンスに問題がある**

## リードモデル

- パフォーマンス低下を回避するために仕様、リポジトリ、ドメインオブジェクトを利用しないことも考える必要が出てくる
- 全てのインスタンスを DB から取得し、マッチするものだけを探すとパフォーマンスが落ちる可能性がある

  - つまり、UX が悪くなり、そのサービスが使われなくなる可能性がある

- そういう時は、ドメイン駆動設計を忠実に守るのではなく、パフォーマンスを優先する必要もある
- **必要なデータだけを検索するサービスオブジェクト的なものを作成することでパフォーマンス問題を回避する**
  - DB を操作する SQL とルールのデータを併用して記述することで必要なデータのみ取得

### CQS(Command Query Separation)

- データの検索は複雑な場合が多いが、ドメインに対する制約が少なく動作も単純なので、ドメインオブジェクトの利用を緩和する
- データの保存はドメインとしての制約も多いので、ドメインオブジェクトを積極的に利用する

## 13 章 まとめ

- ドメインオブジェクト自身のメソッドで表すのに複雑な場合は仕様として抽出する

# 14 章 アーキテクチャー

- 何をどこに記述するのかの設計方針
- レイヤードアーキテクチャ
- ヘキサゴナルアーキテクチャー
- クリーンアーキテクチャー

## レイヤードアーキテクチャー

- 以下の層に分けて表現

1. ユーザインターフェース層(プレゼンテーション層)
2. アプリケーション層
3. ドメイン層
4. インフラストラクチャー

- **依存関係の方向は上から下**
  - 1->2 は OK
  - 4->1 は NG

### ユーザインターフェース(プレゼンテーション層)

- ユーザインターフェースとアプリケーションをつなげる部分
- ユーザインターフェースはなんでも可能
- MVC のコントローラー部分
- クライアント入力の値の解釈と結果表示

### アプリケーション層

- ユースケースの進行役
- 進行役なので、ドメインのルール、ロジックは書かない
- ドメイン層を取りまとめる唯一の存在

### ドメイン層

- ドメインのルールを記述するところ
- ファクトリとリポジトリのインターフェースを持つ

### インフラストラクチャー

- 他の層を支える技術的基盤
- リポジトリ、ファクトリの実装

## ヘキサゴナルアーキテクチャー

- アプリケーション以外は取り外し可能とするアーキテクチャー
- 特定の DB に依存しない
- 特定のユーザインターフェースに依存しない

## クリーンアーキテクチャー

- アプリケーション以外は取り外し可能とするアーキテクチャーの一つ

# 15 章 モデリング

- パターンを取り入れようとしすぎてドメインの本質を見過ごしては本末転倒
- ドメインの本質を見つめるべき

## ドメインエキスパート(企画)の人とモデリング

- 人の会話は知らず知らず認識のずれが起きる

1. 企画部と開発部がシステム開発のための会議をする
2. 企画部の人は、開発部の言葉が難しくて理解できない
3. あまり理解できていないが、システム開発部のエキスパートである開発部に一旦任せようと思う
4. 何事もなかったように質問もされない
5. ドメインが捻じ曲げられたものがコードに反映され、システムができる
6. 見当違いのものができる

- そうならないために以下 2 点を挙げる
  - ユビキタス言語
  - コンテキスト

### ユビキタス言語

- 企画部と開発部などで共通の言葉を決定する
- その言語はドキュメントや会話、コード上全てで共通して利用する必要がある

### コンテキスト

- ドメインの国境となるもので異なるコンテキストでは同じものもニュアンス異なる
- 例）ユーザー
  - 認証という観点
    - id, パスワード
  - 人としての観点
    - サークルを作る
    - 名前を変更する
- これらはそれぞれの名前空間で分けると分かりやすくなる
- しかし、分けすぎるとそれぞれはつながりが分かりにくくなるため、コンテキストマップを作成する必要がある

### 実装

- トップダウンで仕様を落とし込み、ボトムアップで実装していく
