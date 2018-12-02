package main

import (
	"html/template"
	"io/ioutil" // golangで他のテキストファイルを読み込んだり書き込んだりできる
	"net/http"
)

// wikiの構造体
type Page struct {
	Title string //タイトル
	Body  []byte //内容
}

// パスのアドレスを設定して文字の長さを定数として持つ
const lenPath = len("/view/")

// テンプレートファイルの配列を作成
var templates = make(map[string]*template.Template)

//初期化関数
func init() {
	for _, tmpl := range []string{"edit", "view"} {
		// エラーの場合Panicを起こすためのエラー処理はなし
		t := template.Must(template.ParseFiles(tmpl + ".html"))
		templates[tmpl] = t
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, err := loadPage(title)
	if err != nil {
		// editHandlerのURLに飛ばすことで編集ページに飛ばすことができる
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// 編集ページのパス
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	// edit.htmlファイルからのbodyの取得
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		// statusをInternalServerErrorとして出力
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//作成したviewpageにリダイレクト
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// tmpl.html内にTitleやBodyを入れるようにする
	err := templates[tmpl].Execute(w, p)
	if err != nil {
		// statusをInternalServerErrorとして出力
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// テキストファイル保存メソッド
func (p *Page) save() error { //errorが発生したらエラーをだすようにする
	// タイトルの名前でテキストファイルを作成し、保存する(重複をなくす)
	filename := p.Title + ".txt"
	// テキストファイルを作成する
	// 0600はテキストデータを書き込んだり読み込んだり設定する
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// titleからファイル名を読み込んで新しいPageのポインタを返す
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	// errに値が入ったらエラーとしてbodyの値にnilにして返す
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
