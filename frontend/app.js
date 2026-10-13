document.addEventListener("DOMContentLoaded", () => {
    const itemsList = document.getElementById("items-list");
    const statusDot = document.getElementById("status-dot");
    const statusText = document.getElementById("status-text");
    const refreshBtn = document.getElementById("refresh-btn");
    const loader = document.getElementById("loader");
    const emptyState = document.getElementById("empty-state");

    async function fetchItems() {
        loader.classList.remove("hidden");
        emptyState.classList.add("hidden");

        try {
            // Запрашиваем эндпоинт бэкенда
            const response = await fetch("/api/items");
            
            if (!response.ok) {
                throw new Error(`Ошибка HTTP: ${response.status}`);
            }

            const data = await response.json();

            // Статус успешен
            statusDot.className = "dot online";
            statusText.textContent = "Бэкенд и БД в сети";

            renderItems(data);
        } catch (error) {
            console.error("Ошибка при получении данных:", error);
            statusDot.className = "dot offline";
            statusText.textContent = "Ошибка соединения";
            itemsList.innerHTML = "";
            emptyState.textContent = "Не удалось подключиться к серверу или базе данных.";
            emptyState.classList.remove("hidden");
        } finally {
            loader.classList.add("hidden");
        }
    }

    function renderItems(items) {
        itemsList.innerHTML = "";

        if (!items || items.length === 0) {
            emptyState.textContent = "В базе данных пока нет записей.";
            emptyState.classList.remove("hidden");
            return;
        }

        items.forEach((item) => {
            const li = document.createElement("li");
            li.className = "item-row";
            li.innerHTML = `
                <span class="item-name">${escapeHtml(item.name)}</span>
                <span class="item-id">#${item.id}</span>
            `;
            itemsList.appendChild(li);
        });
    }

    // Защита от XSS при выводе текста
    function escapeHtml(str) {
        const div = document.createElement("div");
        div.textContent = str;
        return div.innerHTML;
    }

    refreshBtn.addEventListener("click", fetchItems);

    // Первичная загрузка при старте страницы
    fetchItems();
});