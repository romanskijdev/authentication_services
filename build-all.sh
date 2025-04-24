#!/bin/bash
# 💀 для CRITICAL (критические) уязвимости
# 🔥 для HIGH (высокие) уязвимости
# ⚠️ для MEDIUM (средние) уязвимости и предупреждения
# ℹ️ для LOW (низкие) уязвимости

set -e
trap 'echo "🔴 Скрипт остановлен"; exit' SIGINT

# Определяем корневую директорию скрипта
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SECURITY_REPORTS_DIR="${SCRIPT_DIR}/security_reports"

# Функция для форматирования отчета с смайлами
format_report() {
    local report_file=$1
    local output_file=$2
    
    # Создаем временные файлы для разных уровней критичности
    local critical_file=$(mktemp)
    local high_file=$(mktemp)
    local medium_file=$(mktemp)
    local low_file=$(mktemp)
    local other_file=$(mktemp)
    
    # Обрабатываем каждую строку отчета и распределяем по файлам
    while IFS= read -r line; do
        if [[ $line =~ "CRITICAL" ]]; then
            echo "💀 $line" >> "$critical_file"
        elif [[ $line =~ "HIGH" ]]; then
            echo "🔥 $line" >> "$high_file"
        elif [[ $line =~ "MEDIUM" ]]; then
            echo "⚠️ $line" >> "$medium_file"
        elif [[ $line =~ "LOW" ]]; then
            echo "ℹ️ $line" >> "$low_file"
        elif [[ $line =~ "Error" ]]; then
            echo "❌ $line" >> "$critical_file"
        elif [[ $line =~ "Warning" ]]; then
            echo "⚠️ $line" >> "$medium_file"
        else
            echo "$line" >> "$other_file"
        fi
    done < "$report_file"
    
    # Объединяем файлы в правильном порядке и перезаписываем итоговый файл
    {
        echo "🕒 Время проверки (UTC): $(date -u '+%Y-%m-%d %H:%M:%S')"
        echo "======================================="
        echo
        echo "=== КРИТИЧЕСКИЕ УЯЗВИМОСТИ ==="
        cat "$critical_file"
        echo
        echo "=== ВЫСОКИЕ УЯЗВИМОСТИ ==="
        cat "$high_file"
        echo
        echo "=== СРЕДНИЕ УЯЗВИМОСТИ ==="
        cat "$medium_file"
        echo
        echo "=== НИЗКИЕ УЯЗВИМОСТИ ==="
        cat "$low_file"
        echo
        echo "=== ДРУГАЯ ИНФОРМАЦИЯ ==="
        cat "$other_file"
    } > "$output_file"
    
    # Удаляем временные файлы
    rm -f "$critical_file" "$high_file" "$medium_file" "$low_file" "$other_file"
}

# Функция для сканирования кода с помощью gosec
scan_code() {
    local service_name=$1
    echo "🔍 Сканирование кода $service_name на уязвимости..."
    
    local service_reports_dir="${SECURITY_REPORTS_DIR}/${service_name}"
    mkdir -p "${service_reports_dir}"
    
    # Проверяем, установлен ли gosec
    if ! command -v gosec &> /dev/null; then
        echo "📦 Устанавливаем gosec..."
        go install github.com/securego/gosec/v2/cmd/gosec@latest
    fi
    
    # Сканируем код и сохраняем отчет
    local temp_report="${service_reports_dir}/code_report.tmp"
    gosec -fmt=text -out="$temp_report" ./... || true
    
    # Форматируем отчет с смайлами и сортировкой
    format_report "$temp_report" "${service_reports_dir}/code_report.txt"
    rm -f "$temp_report"
    
    echo "📝 Отчет по коду сохранен в ${service_reports_dir}/code_report.txt"
}

# Функция для сканирования Docker образа
scan_image() {
    local image_name=$1
    local service_name=$2
    echo "🔍 Сканирование образа $image_name на уязвимости..."
    
    local service_reports_dir="${SECURITY_REPORTS_DIR}/${service_name}"
    mkdir -p "${service_reports_dir}"
    
    local temp_report="${service_reports_dir}/image_report.tmp"
    trivy image --severity HIGH,CRITICAL \
        --format table \
        --output "$temp_report" \
        "${image_name}" || true
    
    # Форматируем отчет с смайлами и сортировкой
    format_report "$temp_report" "${service_reports_dir}/image_report.txt"
    rm -f "$temp_report"
    
    echo "📝 Отчет по образу сохранен в ${service_reports_dir}/image_report.txt"
}

# Функция для проверки на наличие секретов в коде
scan_secrets() {
    local service_name=$1
    echo "🔍 Проверка $service_name на наличие секретов в коде..."
    
    local service_reports_dir="${SECURITY_REPORTS_DIR}/${service_name}"
    mkdir -p "${service_reports_dir}"
    
    # Проверяем, установлен ли gitleaks
    if ! command -v gitleaks &> /dev/null; then
        echo "📦 Устанавливаем gitleaks..."
        go install github.com/zricethezav/gitleaks/v8@latest
    fi
    
    # Сканируем код на наличие секретов
    local temp_report="${service_reports_dir}/secrets.tmp"
    
    # Создаем временный файл с исключениями из .gitignore
    local gitignore_excludes=""
    if [ -f ".gitignore" ]; then
        gitignore_excludes=$(grep -v '^#' .gitignore | grep -v '^$' | sed 's/^/--exclude=/')
    fi
    
    # Запускаем gitleaks с учетом .gitignore
    gitleaks detect --source . --report-format json --no-git $gitignore_excludes --report-path "$temp_report" || true
    
    # Форматируем отчет с смайлами
    {
        echo "🕒 Время проверки (UTC): $(date -u '+%Y-%m-%d %H:%M:%S')"
        echo "======================================="
        echo
        echo "=== ПРОВЕРКА НА СЕКРЕТЫ В КОДЕ ==="
        echo
        if [ -s "$temp_report" ]; then
            echo "🔴 ОБНАРУЖЕНЫ СЕКРЕТЫ В КОДЕ!"
            echo "Типы обнаруженных секретов:"
            echo
            # Парсим JSON и добавляем смайлы к разным типам секретов
            jq -r '.[] | select(.File | test("^[^/]+/")) | "\(.RuleID) - \(.Secret)\nФайл: \(.File)\nСтрока: \(.StartLine)\nКонтекст: \(.Context)\n"' "$temp_report" | while IFS= read -r line; do
                if [[ $line =~ "API_KEY" ]]; then
                    echo "🔑 $line"
                elif [[ $line =~ "PASSWORD" ]]; then
                    echo "🔐 $line"
                elif [[ $line =~ "SECRET" ]]; then
                    echo "🤫 $line"
                elif [[ $line =~ "TOKEN" ]]; then
                    echo "🎫 $line"
                elif [[ $line =~ "PRIVATE_KEY" ]]; then
                    echo "🔏 $line"
                elif [[ $line =~ "CREDENTIAL" ]]; then
                    echo "📝 $line"
                elif [[ $line =~ ^Файл: ]]; then
                    echo "📄 $line"
                elif [[ $line =~ ^Строка: ]]; then
                    echo "📏 $line"
                elif [[ $line =~ ^Контекст: ]]; then
                    echo "📋 $line"
                    echo "----------------------------------------"
                else
                    echo "⚠️ $line"
                fi
            done
        else
            echo "🟢 Секреты в коде не обнаружены"
        fi
    } > "${service_reports_dir}/secrets_report.txt"
    
    rm -f "$temp_report"
    echo "📝 Отчет по секретам сохранен в ${service_reports_dir}/secrets_report.txt"
}

# Функция для сборки core
build_core() {
    local service_name=$1
    echo "🧡 $service_name"
    gofmt -w .
    if go mod tidy && go mod vendor && go vet ./...; then
        echo "🟢 $service_name успешно собран"
        # Сканирование кода и секретов для core
        scan_code $service_name
        scan_secrets $service_name
    else
        echo "🔴 Ошибка при сборке $service_name"
        return 1
    fi
}

# Функция для сборки сервисов
build_service() {
    local service_name=$1
    echo "🧡 $service_name"
    gofmt -w .
    if go mod tidy && go mod vendor && go vet ./...; then
        echo "🟢 $service_name успешно собран"
        
        # Сканируем код и секреты
        scan_code $service_name
        scan_secrets $service_name
        
        # Проверяем наличие Dockerfile.serv
        if [ ! -f "Dockerfile" ]; then
            echo "⚠️ Файл Dockerfile не найден в директории $(pwd)"
            return 0
        fi
        
        # Собираем Docker образ с использованием BuildKit
        echo "🔨 Сборка Docker образа..."
        DOCKER_BUILDKIT=1 docker build \
            --progress=plain \
            --build-arg APP_NAME="${service_name}_app" \
            --build-arg BUILDKIT_INLINE_CACHE=1 \
            -t "${service_name}_i" \
            -f Dockerfile .
        
        if [ $? -eq 0 ]; then
            # Сканируем образ на уязвимости
            scan_image "${service_name}_i" "$service_name"
        else
            echo "🔴 Ошибка при сборке Docker образа"
        fi
    else
        echo "🔴 Ошибка при сборке $service_name"
        return 1
    fi
}

# Включаем BuildKit
export DOCKER_BUILDKIT=1

cd core
./build_proto.sh
build_core "core"  # Используем специальную функцию для core

echo "🧡 ls: "
ls

cd ../services

cd admin_service
build_service "admin_service"

cd ../notification_service
build_service "notification_service"

cd ../payment_service
build_service "payment_service"

cd ../system-service
build_service "system-service"

cd ../tariff-service
build_service "tariff-service"

cd ../telegram_service
build_service "telegram_service"

cd ../user-service
build_service "user-service"

cd ../w-imap-service
build_service "w-imap-service"

cd ../w-smtp-service
build_service "w-smtp-service"

echo "🟢 Все сервисы обработаны"

# Генерация общего отчета
echo "📊 Генерация общего отчета..."
{
    echo "=== Отчет о безопасности $(date) ==="
    echo "======================================="
    echo
    echo "=== АНАЛИЗ CORE ==="
    if [ -f "${SECURITY_REPORTS_DIR}/core/code_report.txt" ]; then
        cat "${SECURITY_REPORTS_DIR}/core/code_report.txt"
    fi
    if [ -f "${SECURITY_REPORTS_DIR}/core/secrets_report.txt" ]; then
        cat "${SECURITY_REPORTS_DIR}/core/secrets_report.txt"
    fi
    echo "----------------------------------------"
    
    for service in admin_service notification_service payment_service system-service tariff-service telegram_service user-service w-imap-service w-smtp-service; do
        echo "=== $service ==="
        echo "Анализ кода:"
        if [ -f "${SECURITY_REPORTS_DIR}/${service}/code_report.txt" ]; then
            cat "${SECURITY_REPORTS_DIR}/${service}/code_report.txt"
        fi
        echo "Анализ секретов:"
        if [ -f "${SECURITY_REPORTS_DIR}/${service}/secrets_report.txt" ]; then
            cat "${SECURITY_REPORTS_DIR}/${service}/secrets_report.txt"
        fi
        echo "Анализ Docker образа:"
        if [ -f "${SECURITY_REPORTS_DIR}/${service}/image_report.txt" ]; then
            cat "${SECURITY_REPORTS_DIR}/${service}/image_report.txt"
        fi
        echo "----------------------------------------"
    done
} > "${SECURITY_REPORTS_DIR}/full_security_report.txt"

echo "📝 Полный отчет доступен в ${SECURITY_REPORTS_DIR}/full_security_report.txt"

# Проверка на критические уязвимости и секреты
if grep -q "💀\|🔥\|🔴 ОБНАРУЖЕНЫ СЕКРЕТЫ В КОДЕ" "${SECURITY_REPORTS_DIR}/full_security_report.txt"; then
    echo "⚠️ ВНИМАНИЕ: Обнаружены критические или высокие уязвимости или секреты в коде!"
    echo "Пожалуйста, просмотрите отчет и примите необходимые меры"
fi