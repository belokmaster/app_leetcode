// Функция загрузки задач
async function loadTasks() {
    try {
        const response = await fetch('/api/tasks');
        const tasks = await response.json();

        const container = document.getElementById('tasks-container');
        if (tasks.length === 0) {
            container.innerHTML = '<div class="empty-state">Пока нет задач. Добавьте первую!</div>';
            return;
        }

        container.innerHTML = tasks.map(task => `
            <div class="task-item" id="task-${task.id}">
                <button class="delete-btn" onclick="deleteTask(${task.id})">×</button>
                <div class="task-number">Задача #${task.number}</div>
                <div class="task-desc">${task.description}</div>
                
                <div class="task-meta">
                    <span class="difficulty">Сложность: ${task.platform_difficult}/3</span>
                    <span class="difficulty">Моя: ${task.my_difficult}/10</span>
                    <span class="status ${task.solved_with_hint ? 'solved-hint' : 'solved-alone'}">
                        ${task.solved_with_hint ? 'С подсказкой' : 'Самостоятельно'}
                    </span>
                    <span class="status ${task.is_masthaved ? 'mastered' : 'not-mastered'}">
                        ${task.is_masthaved ? 'Освоено' : 'Не освоено'}
                    </span>
                </div>
                
                <div class="task-date">Создано: ${new Date(task.created_at).toLocaleDateString('ru-RU')}</div>
            </div>
        `).join('');
    } catch (error) {
        document.getElementById('tasks-container').innerHTML = '<div class="empty-state">Ошибка загрузки задач</div>';
    }
}

// Функция удаления задачи
async function deleteTask(id) {
    if (!confirm('Удалить задачу?')) {
        return;
    }

    try {
        const formData = new URLSearchParams();
        formData.append('id', id.toString());

        const response = await fetch('/tasks/delete', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData
        });

        if (response.ok) {
            const taskElement = document.getElementById(`task-${id}`);
            if (taskElement) {
                taskElement.style.opacity = '0';
                setTimeout(() => taskElement.remove(), 300);
            }
        } else {
            alert('Ошибка удаления');
        }
    } catch (error) {
        alert('Ошибка сети');
    }
}

// Загружаем задачи при загрузке страницы
document.addEventListener('DOMContentLoaded', function () {
    loadTasks();

    // Обновляем каждые 10 секунд
    setInterval(loadTasks, 10000);

    // Обновляем после добавления новой задачи
    document.querySelector('form').addEventListener('submit', function () {
        setTimeout(loadTasks, 1000);
    });
});