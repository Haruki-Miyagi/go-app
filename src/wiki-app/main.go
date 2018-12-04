package main

import (
	"errors"
	"html/template"
	"io/ioutil" // golangで他のテキストファイルを読み込んだり書き込んだりできる
	"log"
	"net/http"
	"regexp" // 文字列の正規表現
	"strings"
)

// wikiの構造体
type Page struct {
	Title string //タイトル
	Body  []byte //内容
}

// パスのアドレスを設定して文字の長さを定数として持つ
const lenPath = len("/view/")

// .txt
const expend_string = ".txt"

// テンプレートファイルの配列を作成
var templates = make(map[string]*template.Template)

// 正規表現でURLを生成できる大文字小文字の英字と数字を判別する
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

// 初期化関数
func init() {
	for _, tmpl := range []string{"edit", "view"} {
		// エラーの場合Panicを起こすためのエラー処理はなし
		t := template.Must(template.ParseFiles(tmpl + ".html"))
		templates[tmpl] = t
	}
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	// main.goがいる階層のディレクトリにある.txtデータを取得する
	files, err := ioutil.ReadDir("./")
	if err != nil {
		err = errors.New("所定のディレクトリ内にテキストファイルがありません")
		log.Print(err)
		return
	}

	var paths []string    // テキストデータの名前
	var fileName []string // テキストデータのファイル名
	for _, file := range files {
		// 対象となる.txtデータのみ取得する
		if strings.HasSuffix(file.Name(), expend_string) {
			fileName = strings.Split(string(file.Name()), expend_string)
			paths = append(paths, fileName[0])
		}
	}

	if paths == nil {
		err = errors.New("テキストファイルが存在しません")
		log.Print(err)
	}

	t := template.Must(template.ParseFiles("top.html"))
	err = t.Execute(w, paths)
	if err != nil {
		// statusをInternalServerErrorとして出力
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		// editHandlerのURLに飛ばすことで編集ページに飛ばすことができる
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// 編集ページのパス
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// edit.htmlファイルからのbodyの取得
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		// statusをInternalServerErrorとして出力
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 作成したviewpageにリダイレクト
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// タイトルのチェックを行う(バリデーション)
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Requestからページタイトルを取り出して、fnを呼び出す
		title := r.URL.Path[lenPath:]
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			err := errors.New("Invalid Page Title")
			log.Print(err)
			return
		}
		fn(w, r, title)
	}
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
	http.HandleFunc("/top/", topHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
