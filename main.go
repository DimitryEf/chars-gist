/*
Исходные данные:

Есть сервер с «многоядерным» процессором
Есть папка, которая содержит неопределённое, но очень больше количество текстовых файлов.
Внутри ASCII символы, пробелы и переводы строк

Нужно: Построить «гистограмму» распределения ASCII символов в этих файлах, т.е. посчитать сколько раз каждый символ встречается в файлах.

При решении задачи временем чтения файлов с диска можно пренебречь, т.е. считать что файлы считываются с диска в память "мгновенно"
*/
package main

import "chars-gist/cmd"

func main() {
	cmd.Execute()
}