const labelMap = {
    0: "Massive", 1: "String", 2: "HashTable", 3: "Math", 4: "DynamicProgramming",
    5: "Sorting", 6: "Greedy", 7: "DepthFirstSearch", 8: "BinarySearch", 9: "DataBase",
    10: "Matrix", 11: "BitManipulation", 12: "Tree", 13: "BreadthFirstSearch",
    14: "TwoPointers", 15: "PrefixSum", 16: "Heap", 17: "Simulation", 18: "Counting",
    19: "Graph", 20: "BinaryTree", 21: "Stack", 22: "SlidingWindow", 23: "Design",
    24: "Enumeration", 25: "Backtracking", 26: "UnionFind", 27: "NumberTheory",
    28: "LinkedList", 29: "OrderedSet", 30: "SegmentTree", 31: "MonotonicStack",
    32: "Trie", 33: "DivideAndConquer", 34: "Combinatorics", 35: "Bitmask",
    36: "Queue", 37: "Recursion", 38: "Geometry", 39: "BinaryIndexedTree",
    40: "Memoization", 41: "HashFunction", 42: "BinarySearchTree", 43: "ShortestPath",
    44: "StringMatching", 45: "TopologicalSort", 46: "RollingHash", 47: "GameTheory",
    48: "Interactive", 49: "DataStream", 50: "MonotonicQueue", 51: "Brainteaser",
    52: "DoubleLinkedList", 53: "MergeSort", 54: "Randomized", 55: "CountingSort",
    56: "Iterator", 57: "Concurrency", 58: "SuffixArray", 59: "LineSweep",
    60: "ProbabilityAndStatistics", 61: "Quickselect", 62: "MinimumSpanningTree",
    63: "BucketSort", 64: "Shell"
};

function populateLabelSelects() {
    const selects = [
        document.getElementById('add-labels'),
        document.getElementById('edit-labels')
    ].filter(Boolean);

    const optionsHTML = Object.entries(labelMap)
        .map(([value, text]) => `<option value="${value}">${text}</option>`)
        .join('');

    selects.forEach((select) => {
        select.innerHTML = optionsHTML;
    });
}

function updateStats(tasks) {
    const total = tasks.length;
    const solved = tasks.filter(task => Boolean(task.solved_at)).length;
    const alone = tasks.filter(task => !task.solved_with_hint).length;
    const mastered = tasks.filter(task => task.is_masthaved).length;

    document.getElementById('stat-total').textContent = total;
    document.getElementById('stat-solved').textContent = solved;
    document.getElementById('stat-alone').textContent = alone;
    document.getElementById('stat-mastered').textContent = mastered;
}

function showMessage(text, isError = false) {
    const toast = document.createElement('div');
    toast.textContent = text;
    toast.style.position = 'fixed';
    toast.style.right = '16px';
    toast.style.bottom = '16px';
    toast.style.padding = '10px 14px';
    toast.style.borderRadius = '10px';
    toast.style.fontWeight = '700';
    toast.style.zIndex = '1200';
    toast.style.background = isError ? '#fce8ea' : '#dff4ef';
    toast.style.color = isError ? '#a83944' : '#0f756d';
    toast.style.border = `1px solid ${isError ? '#e8aab2' : '#95d5c9'}`;
    toast.style.boxShadow = '0 10px 24px rgba(15, 27, 38, 0.15)';

    document.body.appendChild(toast);
    setTimeout(() => toast.remove(), 2200);
}


function getLabelsHTML(labels) {
    if (!labels || labels.length === 0) {
        return '';
    }
    const labelsHTML = labels.map(labelId => {
        const labelName = labelMap[labelId] || `Unknown Label #${labelId}`;
        return `<span class="task-label">${labelName}</span>`;
    }).join('');
    return `<div class="task-labels">${labelsHTML}</div>`;
}


// Функция загрузки задач
async function loadTasks() {
    try {
        const response = await fetch('/api/tasks');
        const tasks = await response.json();
        updateStats(tasks || []);

        const container = document.getElementById('tasks-container');
        if (!tasks || tasks.length === 0) {
            container.innerHTML = '<div class="empty-state">Пока нет задач. Добавьте первую!</div>';
            return;
        }

        container.innerHTML = tasks.map(task => `
            <div class="task-item ${getDifficultyClass(task.platform_difficult)}" id="task-${task.id}" onclick='openEditModal(${JSON.stringify(task).replace(/'/g, "&apos;").replace(/"/g, "&quot;")})'>
                <button class="delete-btn" onclick="event.stopPropagation(); deleteTask(${task.id})">×</button>
                <div class="task-number">Задача #${task.number}</div>
                ${getLabelsHTML(task.labels)}
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
        updateStats([]);
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
            showMessage('Задача удалена');
            setTimeout(loadTasks, 320);
        } else {
            showMessage('Ошибка удаления', true);
        }
    } catch (error) {
        showMessage('Ошибка сети', true);
    }
}

// Загружаем задачи при загрузке страницы
document.addEventListener('DOMContentLoaded', function () {
    populateLabelSelects();
    loadTasks();

    // Обновляем каждые 10 секунд
    setInterval(loadTasks, 10000);

    // Обновляем после добавления новой задачи
    document.getElementById('add-task-form').addEventListener('submit', function () {
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
                    `<div class="empty-state">Задача с номером ${number} не найдена.</div>`;
            } else {
                throw new Error(`Ошибка сервера: ${response.status}`);
            }

            return;
        }

        const data = await response.json();

        resultsContainer.innerHTML = `
            <div class="search-result-item task-item ${getDifficultyClass(data.platform_difficult)}" onclick='openEditModal(${JSON.stringify(data).replace(/'/g, "&apos;").replace(/"/g, "&quot;")})'>
                <div class="task-number">Задача #${data.number}</div>
                ${getLabelsHTML(data.labels)}
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
            `<div class="empty-state">Не удалось выполнить поиск. Проверьте соединение с сервером.</div>`;
    }
}

// Открытие модального окна при клике на задачу
function openEditModal(task) {
    const modal = document.getElementById('edit-modal');

    // Заполняем форму данными задачи
    document.getElementById('edit-id').value = task.id;
    document.getElementById('edit-number').value = task.number;
    document.getElementById('edit-platform-difficult').value = task.platform_difficult;
    document.getElementById('edit-my-difficult').value = task.my_difficult;
    document.getElementById('edit-description').value = task.description;
    document.getElementById('edit-solved-with-hint').checked = task.solved_with_hint;
    document.getElementById('edit-is-masthaved').checked = task.is_masthaved;

    const labelsSelect = document.getElementById('edit-labels');
    const taskLabels = task.labels || [];
    for (let i = 0; i < labelsSelect.options.length; i++) {
        const option = labelsSelect.options[i];
        option.selected = taskLabels.includes(parseInt(option.value));
    }


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

    const labelsSelect = document.getElementById('edit-labels');
    const selectedLabels = Array.from(labelsSelect.selectedOptions).map(option => option.value);

    formData.append('id', document.getElementById('edit-id').value);
    formData.append('platform_difficult', document.getElementById('edit-platform-difficult').value);
    formData.append('my_difficult', document.getElementById('edit-my-difficult').value);
    formData.append('description', document.getElementById('edit-description').value);
    formData.append('solved_at', document.getElementById('edit-solved-at').value);
    formData.append('solved_with_hint', document.getElementById('edit-solved-with-hint').checked ? 'on' : '');
    formData.append('is_masthaved', document.getElementById('edit-is-masthaved').checked ? 'on' : '');
    formData.append('labels', selectedLabels.join(','));


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
            showMessage('Задача обновлена!');
        } else {
            showMessage('Ошибка обновления задачи', true);
        }
    } catch (error) {
        showMessage('Ошибка сети', true);
    }
}

async function getRandomOldTask() {
    try {
        const response = await fetch('/api/tasks/random-old');
        const data = await response.json();

        const resultContainer = document.getElementById('random-task-result');

        if (data.error) {
            resultContainer.innerHTML = `<div class="empty-state">${data.error}</div>`;
            return;
        }

        resultContainer.innerHTML = `
            <div class="random-task-item task-item ${getDifficultyClass(data.platform_difficult)}" onclick='openEditModal(${JSON.stringify(data).replace(/'/g, "&apos;").replace(/"/g, "&quot;")})'>
                <div class="task-number">Задача #${data.number}</div>
                ${getLabelsHTML(data.labels)}
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
            `<div class="empty-state">Ошибка загрузки</div>`;
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