const tg = window.Telegram.WebApp;
tg.expand();

// Инициализация приложения
document.addEventListener('DOMContentLoaded', () => {
    initNavigation();
    initSwipeHandler();
    loadTasks();
});

// Навигация между страницами
function initNavigation() {
    const navButtons = document.querySelectorAll('.nav-btn');
    navButtons.forEach(button => {
        button.addEventListener('click', () => {
            const targetPage = button.dataset.page;
            switchPage(targetPage);
        });
    });
}

function switchPage(pageId) {
    const pages = document.querySelectorAll('.page');
    const navButtons = document.querySelectorAll('.nav-btn');
    
    pages.forEach(page => {
        page.classList.remove('active');
    });
    
    navButtons.forEach(btn => {
        btn.classList.remove('active');
    });
    
    document.getElementById(pageId).classList.add('active');
    document.querySelector(`[data-page="${pageId}"]`).classList.add('active');
}

// Обработка свайпов
function initSwipeHandler() {
    let startY = 0;
    let currentCard = null;

    document.querySelector('.task-container').addEventListener('touchstart', (e) => {
        startY = e.touches[0].clientY;
        currentCard = e.target.closest('.task-card');
    });

    document.querySelector('.task-container').addEventListener('touchmove', (e) => {
        if (!currentCard) return;
        
        const deltaY = e.touches[0].clientY - startY;
        if (deltaY > 0) return; // Запрещаем свайп вверх
        
        currentCard.style.transform = `translateY(${deltaY}px)`;
    });

    document.querySelector('.task-container').addEventListener('touchend', (e) => {
        if (!currentCard) return;
        
        const deltaY = e.changedTouches[0].clientY - startY;
        if (deltaY < -100) {
            // Свайп вниз - показываем следующую карточку
            currentCard.style.transform = 'translateY(-100vh)';
            setTimeout(() => {
                currentCard.remove();
                loadNextTask();
            }, 300);
        } else {
            // Возвращаем карточку на место
            currentCard.style.transform = '';
        }
    });
}

// Загрузка заданий
function loadTasks() {
    fetch('/api/tasks')
        .then(response => response.json())
        .then(tasks => {
            tasks.forEach(task => createTaskCard(task));
        });
}

function createTaskCard(task) {
    const card = document.createElement('div');
    card.className = 'task-card';
    card.innerHTML = `
        <h2>${task.title}</h2>
        <p>${task.description}</p>
        <div class="options">
            ${task.options.map(option => `
                <button class="option-btn" data-correct="${option.isCorrect}">
                    ${option.text}
                </button>
            `).join('')}
        </div>
    `;
    
    document.querySelector('.task-container').appendChild(card);
}

function loadNextTask() {
    fetch('/api/next-task')
        .then(response => response.json())
        .then(task => {
            if (task) {
                createTaskCard(task);
            }
        });
} 