#Использовать messenger

Мессенджер = Новый Мессенджер();
Мессенджер.ИнициализироватьТранспорт("telegram", Новый Структура("Логин", АргументыКоманднойСтроки[0]));
Мессенджер.ОтправитьСообщение("telegram", АргументыКоманднойСтроки[1], АргументыКоманднойСтроки[2], , "html");