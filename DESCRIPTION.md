Bu ödevde SOAP ve/veya REST web servis tabanlı bir e-barter sistemi tasarlamanız ve uygulamanız istenmektedir. Ekte verilen örnek sistemi temel olarak alabilirsiniz.  Öncelikli olarak, Barter para kullanmaksızın ürünlerin değişimiyle yapılan ticaret sistemidir. Sistemi kullanabilmek için, kişi veya şirketlerin belirledikleri benzersiz bir kullanıcı adı ile kaydolmaları gerekmektir. Kayıt işlevini gerçekleştirecek olan web servis yeteneği kullanıcıya yine benzersiz ve güvenlikli bir kod gönderecektir. Sistemde herhangi bir işlem gerçekleştirmek isteyen kullanıcı, herhangi bir servis yeteneğini kullanırken bu kullanıcı adı ve kodu da gönderecek ve ancak kod ve kullanıcı adı eşleştiği ve doğrulandığı zaman servis yeteneği kullanılabilecektir. (API Key Mantığı) Sisteme kaydolan kullanıcılar ihtiyaç duydukları malları ve miktarlarını ve benzer şekilde takas için vermeyi planladıkları malları ve miktarlarını sisteme web servis yoluyla bildirecektir. Tüm kayıtlar e-barter servisinde toplanacak ve sisteme dahil olan her bir yeni ihtiyaç veya takas verisinde, sistem daha önceki bekleyen takas talebi bilgilerini kontrol ederek uygun bir eşleşme olup olmadığı bakacaktır. Takas için uygun bir eşleşme bulunursa, sistem iki kullanıcı buluşturup takası gerçekleştirdiği an iki tarafa da bildirim gönderecektir. Eğer uygun bir eşleşme bulunmazsa, sistem bu kullanıcı ve verilerini kayıt defterine işleyecektir.

Bu sistemde en önemli noktalardan biri e-barter servisi iki kullanıcıyı adil bir şekilde eşleştirebilmek adına çeşitli ürün sınıflarına ait envanter servislerini kullanarak, takasta kullanılmaya aday ürünlerin güncel fiyat bilgileri veya belirli edere kaç adet üründen alınabileceği gibi bilgilere ulaşacaktır.   Örnek uygulamada her bir ürün için sabit bir fiyat değeri varken ödev de ise her bir ürünün birimi için bir fiyat aralığı verilmesi beklenmektedir.  Burada e-barter sisteminden asgari olarak kullanıcı kayıttı alması ve istenen ve verilen mallardan tam olarak gerçekleşebilecek bir eşleşme bulunduğunda bunu gerçekleştirmesidir. Örneğin, A kişi 20 X isterken 60 Y verebiliyor, B kişide 10 Y isterken 30 X verebiliyor ve Envanter sistemlerinde gelen verilere göre 10 Y nin fiyat aralığı ile 20 X’in fiyat aralığı uyuşuyorsa  bu takası gerçekleştir ve iki takas talep kayıttını da bekleyen takaslardan çıkar. Aşağıda bir hesaplama örneği verilmiştir.

Y malı min değer 10,  Y malı maks değer 15 => 10 Y Mal fiyat aralığı [100,150]

X malı min değer 4,  X malı maks değer 6  => 20 X Mal fiyat aralığı [80,120]

Aralıklar arasında çakışma olduğu için bu barter gerçekleştrielebilir.

Bir ödevin değerlendiremeye girebilecek olması için ödev içerisinde en az 4 envanter servisi ve 1 e-barter servisi bulunmalıdır. Ödevin değerlendirilmesinde, yukarıda ifade edilen yeteneklere ve asgari servis sayısına sahip bir sistem değerlendirmeye alınarak alt puandan başlanarak (Çalışan ve asgari koşuları sağlayan bir ödev 70’den başlanarak notlandırılacaktır.) aşağıdaki unusurlar göz önünde bulundurularak puanlandırılacaktır.

    Daha fazla envanter servisinin varlığı
    Envanter sistemlerinde bazılarının aslında birden çok alt-envanter sisteminden oluşması (Örneğin; Gıda Envanter sisteminin, Meyve, Kuru-gıda ve Unlu Mamüler Envanter sistemlerinin kompozisyonundan oluşması.)
    Utility servis kullanımı (kullanıcılara bildirim gitmesi için bir servis kullanımı, burada kullanıcıdan kayıt sırasında istenecek bir email bilgisine e-mail atacak veya bu durumu simüle edecek bir servis olabilir)
    E-barter servisinde kısmı ve / veya çoklu takas imkanı (Bir kişi tarafından talep edilen malın müsait miktarının o anki uygun takas bekleyen kişilerden karşılanırken, kalan kısım sonrasında veya o anda başka bir takas bekleyenden karşılanması)
    Servis implementasyonlarının en azından iki platformda gerçekleştirilmesi,
    Ödev için implemente edilen web servislerin taşınabilir olması (projenin sizin bilgisayarınız dan başka bir bilgisayarda çalıştırılabilir olması)

Ödevler bireyseldir ve grup halinde yapılması beklenmemektedir. Bu doğrultuda bir izlenim oluşması halinde ilgili ödevlerden puan kırılacaktır. 

Ödev’in çıktısı olarak;

    Ödevin tüm kaynak kodlarında, projenin kolayca anlaşılmasına yetecek kadar yorum satırı olmalıdır.
    Ödevin tüm kaynak kodları düzenli halde tek bir sıkıştırılmış klasör olarak teslim edilecektir.
    Ödevin neticesinde ortaya çıkan sistemin çalışma anında çekilmiş ve sistemin çalışmasını gösteren bir video çekilerek bu video da sıkıştırılmış klasöre dahil edilecektir. Video 1 dk’dan uzun olmayacaktır.

Başarılar dilerim.