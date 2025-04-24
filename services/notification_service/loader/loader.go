package loader

import (
	typesm "authentication_service/notification_service/types"
	"fmt"
	"html/template"
	"log"
	"os"
)

func LoadMailTemplates() *typesm.TemplatesMailSystem {
	templatesMailObj := &typesm.TemplatesMailSystem{}
	basePath := "loader/mail-template"
	mailTemplatesNameMap := map[string]string{
		"NewDeviceInfoTemplate": "new-device-info.html",
	}

	for key, value := range mailTemplatesNameMap {
		pathTemplate := fmt.Sprintf("%s/%s", basePath, value)
		if _, err := os.Stat(pathTemplate); err != nil { // Проверка существования файла
			log.Println("🔴 Failed to find mail template:", key, err)
			continue // Пропускаем шаблон, если не найден
		}
		fileBytes, err := os.ReadFile(pathTemplate) // Используем os.ReadFile
		if err != nil {
			log.Println("🔴 Failed to load mail template:", key, err)
			continue // Пропускаем шаблон, если не удалось загрузить
		}
		t, err := template.New(key).Parse(string(fileBytes))
		if err != nil {
			log.Println("🔴 Failed to load mail template:", key, err)
			continue // Пропускаем шаблон, если не удалось разобрать
		}
		switch key {
		case "NewDeviceInfoTemplate":
			templatesMailObj.NewDeviceInfoTemplate = t
		}
	}
	return templatesMailObj
}
