package main

import (
	"fmt"
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
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
	http.ListenAndServe(":8080", nil)
}
