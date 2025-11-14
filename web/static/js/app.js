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
            <div class="task-item ${getDifficultyClass(task.platform_difficult)}" id="task-${task.id}" onclick="openEditModal(${JSON.stringify(task).replace(/"/g, '&quot;')})">
                <button class="delete-btn" onclick="event.stopPropagation(); deleteTask(${task.id})">×</button>
                <div class="task-number">Задача #${task.number}</div>
                <div class="task-desc">${task.description}</div>
                
                <div class="task-meta">
                    <span class="difficulty">Сложность: ${task.platform_difficult}/3</span>
                    <span class="difficulty">Моя: ${task.my_difficult}/10</span>
                    <span class="status ${task.solved_with_hint ? 'solved-hint' : 'solved-alone'}">
                        ${task.solved_with_hint ? 'С подсказкой' : 'Самостоятельно'}
                    </span>
                    <span class="status ${task.is_masthaved ? 'mastered' : 'not-mastered'}">
                        ${task.is_masthaved ? 'Нужная задача' : 'Ненужная задача'}
                    </span>
                </div>
                
                <div class="task-date">
                    Создано: ${new Date(task.created_at).toLocaleDateString('ru-RU')}
                    ${task.solved_at ? ` | Решено: ${new Date(task.solved_at).toLocaleDateString('ru-RU')}` : ''}
                </div>
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

    // Исправленный обработчик поиска
    document.getElementById('search-form').addEventListener('submit', function (e) {
        e.preventDefault();
        const numberInput = this.querySelector('input[name="search_number"]');
        const number = parseInt(numberInput.value);
        if (number > 0) {
            searchTaskByNumber(number);
        }
    });

    // Закрытие модального окна
    document.querySelector('.close').addEventListener('click', closeEditModal);
    document.getElementById('edit-form').addEventListener('submit', handleEditFormSubmit);

    // Закрытие по клику вне модального окна
    window.addEventListener('click', function (e) {
        const modal = document.getElementById('edit-modal');
        if (e.target === modal) {
            closeEditModal();
        }
    });

    document.addEventListener('keydown', function (e) {
        const modal = document.getElementById('edit-modal');
        if (e.key === 'Escape' && modal.style.display === 'block') {
            closeEditModal();
        }
    });
});

// Функция поиска задачи по номеру
async function searchTaskByNumber(number) {
    const resultsContainer = document.getElementById('search-results');
    resultsContainer.innerHTML = '';

    try {
        const response = await fetch(`/api/task?number=${number}`);

        if (!response.ok) {
            if (response.status === 404) {
                resultsContainer.innerHTML =
                    `<div style="color: #7f8c8d; text-align: center;">Задача с номером ${number} не найдена.</div>`;
            } else {
                throw new Error(`Ошибка сервера: ${response.status}`);
            }

            return;
        }

        const data = await response.json();

        resultsContainer.innerHTML = `
            <div class="search-result-item task-item ${getDifficultyClass(data.platform_difficult)}" onclick="openEditModal(${JSON.stringify(data).replace(/"/g, '&quot;')})">
                <div class="task-number">Задача #${data.number}</div>
                <div class="task-desc">${data.description}</div>
                <div class="task-meta">
                    <span class="difficulty">Сложность: ${data.platform_difficult}/3</span>
                    <span class="difficulty">Моя: ${data.my_difficult}/10</span>
                    <span class="status ${data.solved_with_hint ? 'solved-hint' : 'solved-alone'}">
                        ${data.solved_with_hint ? 'С подсказкой' : 'Самостоятельно'}
                    </span>
                    <span class="status ${data.is_masthaved ? 'mastered' : 'not-mastered'}">
                        ${data.is_masthaved ? 'Нужная задача' : 'Ненужная задача'}
                    </span>
                </div>
                <div class="task-date">
                    Создано: ${new Date(data.created_at).toLocaleDateString('ru-RU')}
                    ${data.solved_at ? ` | Решено: ${new Date(data.solved_at).toLocaleDateString('ru-RU')}` : ''}
                </div>
            </div>
        `;
    } catch (error) {
        resultsContainer.innerHTML =
            `<div style="color: #e74c3c; text-align: center;">Не удалось выполнить поиск. Проверьте соединение с сервером.</div>`;
    }
}

// Открытие модального окна при клике на задачу
function openEditModal(task) {
    const modal = document.getElementById('edit-modal');
    const form = document.getElementById('edit-form');

    // Заполняем форму данными задачи
    document.getElementById('edit-id').value = task.id;
    document.getElementById('edit-number').value = task.number;
    document.getElementById('edit-platform-difficult').value = task.platform_difficult;
    document.getElementById('edit-my-difficult').value = task.my_difficult;
    document.getElementById('edit-description').value = task.description;
    document.getElementById('edit-solved-with-hint').checked = task.solved_with_hint;
    document.getElementById('edit-is-masthaved').checked = task.is_masthaved;

    if (task.solved_at) {
        const solvedDate = new Date(task.solved_at);
        document.getElementById('edit-solved-at').value = solvedDate.toISOString().split('T')[0];
    } else {
        document.getElementById('edit-solved-at').value = '';
    }

    modal.style.display = 'block';
}

// Закрытие модального окна
function closeEditModal() {
    document.getElementById('edit-modal').style.display = 'none';
}

// Обработчик отправки формы редактирования
async function handleEditFormSubmit(e) {
    e.preventDefault();

    const formData = new URLSearchParams();
    formData.append('id', document.getElementById('edit-id').value);
    formData.append('platform_difficult', document.getElementById('edit-platform-difficult').value);
    formData.append('my_difficult', document.getElementById('edit-my-difficult').value);
    formData.append('description', document.getElementById('edit-description').value);
    formData.append('solved_at', document.getElementById('edit-solved-at').value);
    formData.append('solved_with_hint', document.getElementById('edit-solved-with-hint').checked ? 'on' : '');
    formData.append('is_masthaved', document.getElementById('edit-is-masthaved').checked ? 'on' : '');

    try {
        const response = await fetch('/api/tasks/update', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData
        });

        if (response.ok) {
            closeEditModal();
            loadTasks();
            alert('Задача обновлена!');
            location.reload()
        } else {
            alert('Ошибка обновления задачи');
        }
    } catch (error) {
        alert('Ошибка сети');
    }
}

async function getRandomOldTask() {
    try {
        const response = await fetch('/api/tasks/random-old');
        const data = await response.json();

        const resultContainer = document.getElementById('random-task-result');

        if (data.error) {
            resultContainer.innerHTML = `<div style="color: #e74c3c; text-align: center;">${data.error}</div>`;
            return;
        }

        resultContainer.innerHTML = `
            <div class="random-task-item task-item ${getDifficultyClass(data.platform_difficult)}" onclick="openEditModal(${JSON.stringify(data).replace(/"/g, '&quot;')})">
                <div class="task-number">Задача #${data.number}</div>
                <div class="task-desc">${data.description}</div>
                <div class="task-meta">
                    <span class="difficulty">Сложность: ${data.platform_difficult}/3</span>
                    <span class="difficulty">Моя: ${data.my_difficult}/10</span>
                    <span class="status ${data.solved_with_hint ? 'solved-hint' : 'solved-alone'}">
                        ${data.solved_with_hint ? 'С подсказкой' : 'Самостоятельно'}
                    </span>
                    <span class="status ${data.is_masthaved ? 'mastered' : 'not-mastered'}">
                        ${data.is_masthaved ? 'Нужная задача' : 'Ненужная задача'}
                    </span>
                </div>
                <div class="task-date">
                    Создано: ${new Date(data.created_at).toLocaleDateString('ru-RU')}
                    ${data.solved_at ? ` | Решено: ${new Date(data.solved_at).toLocaleDateString('ru-RU')}` : ''}
                </div>
            </div>
        `;
    } catch (error) {
        document.getElementById('random-task-result').innerHTML =
            `<div style="color: #e74c3c; text-align: center;">Ошибка загрузки</div>`;
    }
}

function getDifficultyClass(difficulty) {
    switch (difficulty) {
        case 1: return 'easy';
        case 2: return 'medium';
        case 3: return 'hard';
        default: return '';
    }
}