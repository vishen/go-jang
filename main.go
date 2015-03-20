package main

func main() {
	url := "http://google.com"
	root, _ := Get(url)
	WalkNode(root)
}
