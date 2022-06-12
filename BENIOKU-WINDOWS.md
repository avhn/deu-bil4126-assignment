## WINDOWS
> go ve gcc PATH'de olmalı.
https://go.dev/dl/ sitesinden golang'i indirin. gcc yoksa MinGW vb. kaynaklardan bunu da edinin. Yüklerken yönetici izni vermeyi unutmayın. (Binary PATH değişkenine kaydedilmeli.)
---
> servisleri çalıştırın.
Adım 1: main.go dosyasının olduğu konumda PowerShell çalıştırılır. Veya cd komutu proje dosyasının içine gelinir.
Adım 2: PowerShell içinde "go get" komutu gönderilir.
(sqlite3 library eski olduğundan warning verebilir, bu bir error değildir.)
Adım 3. PowerShell içinde "go run main.go" komutu gönderilip servisler başlatılır.
---
> Dokümentasyon proje içindeki README.md dosyasındadır. Her endpoint dokümente edilmiştir. Döndürülen status kodlarına önem verildi.
POSTMANde test için olan hazır koleksiyona aşağıdaki public linkten ulaşabilirsiniz, postman içindeki import->link kısmından import edilir:
https://www.getpostman.com/collections/ef6a711c6dc26ac03197

Not: Input raw JSON olarak body'de kabul ediliyor. Son JSON objectin sonunda virgül olursa bunu notasyona aykırı saydığından kabul etmiyor.
{
    "item": a,
    "item2": b
}
örneğindeki gibi b'den sonra virgül olmadığı gibi son itemden sonra virgül koyulmamalı.
