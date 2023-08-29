#Использовать messenger

ЧтениеJSON = Новый ЧтениеJSON();
ЧтениеJSON.ОткрытьФайл("scriptcfg.json");
НастройкиПриложения  = ПрочитатьJSON(ЧтениеJSON);
ЧтениеJSON.Закрыть();

ЧтениеJSON.УстановитьСтроку(АргументыКоманднойСтроки[0]);
ПараметрыУведомления = ПрочитатьJSON(ЧтениеJSON);
ЧтениеJSON.Закрыть();

Если ПараметрыУведомления.Количество() > 0 и ЗначениеЗаполнено(НастройкиПриложения.Группа) и ЗначениеЗаполнено(НастройкиПриложения.Логин) Тогда

    МассивСтрок = Новый Массив;
    Для Каждого СтрокаУведомления из ПараметрыУведомления Цикл
        МассивСтрок.Добавить(
            стрШаблон("Пользователь %1 отпустил в конфигурации %2 объекты %3",
            СтрокаУведомления.user, СтрокаУведомления.configuration, стрСоединить(СтрокаУведомления.objects, ", "))
        );
    КонецЦикла;
    Уведомление = СтрСоединить(МассивСтрок, Символы.ПС);

    Мессенджер = Новый Мессенджер();
    Мессенджер.ИнициализироватьТранспорт("telegram", Новый Структура("Логин", НастройкиПриложения.Логин));
    Мессенджер.ОтправитьСообщение("telegram", НастройкиПриложения.Группа, Уведомление, , "html");

КонецЕсли;