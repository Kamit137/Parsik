document.addEventListener('DOMContentLoaded', function() {
    const loginBtn = document.getElementById('loginBtn');
    const registerBtn = document.getElementById('registerBtn');
    const loginModal = document.getElementById('loginModal');
    const registerModal = document.getElementById('registerModal');
    const closeLogin = document.getElementById('closeLogin');
    const closeRegister = document.getElementById('closeRegister');
    const authButtons = document.getElementById('authButtons');
    const userContainer = document.getElementById('userContainer');
    const userAvatar = document.getElementById('userAvatar');
    const userMenu = document.getElementById('userMenu');
    const logoutBtn = document.getElementById('logoutBtn');
    const profileBtn = document.getElementById('profileBtn');
    
    let currentUser = null;

    // Открытие модального окна входа
    loginBtn.addEventListener('click', function(e) {
        e.preventDefault();
        loginModal.style.display = 'flex';
    });
    
    // Открытие модального окна регистрации
    registerBtn.addEventListener('click', function(e) {
        e.preventDefault();
        registerModal.style.display = 'flex';
    });
    
    // Закрытие модального окна входа
    closeLogin.addEventListener('click', function() {
        loginModal.style.display = 'none';
    });
    
    // Закрытие модального окна регистрации
    closeRegister.addEventListener('click', function() {
        registerModal.style.display = 'none';
    });
    
    // Закрытие при клике вне модального окна
    window.addEventListener('click', function(e) {
        if (e.target === loginModal) {
            loginModal.style.display = 'none';
        }
        if (e.target === registerModal) {
            registerModal.style.display = 'none';
        }
    });
    
    // Обработка формы входа
    document.getElementById('loginForm').addEventListener('submit', function(e) {
        e.preventDefault();
        const email = document.getElementById('loginEmail').value;
        const password = document.getElementById('loginPassword').value;
        
        // Здесь должна быть логика авторизации
        // Для примера просто сохраняем имя пользователя
        currentUser = {
            name: email.split('@')[0], // Берем часть до @ как имя
            email: email
        };
        
        loginModal.style.display = 'none';
        showUserAvatar();
    });
    
    // Обработка формы регистрации
    document.getElementById('registerForm').addEventListener('submit', function(e) {
        e.preventDefault();
        const name = document.getElementById('registerName').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;
        const confirmPassword = document.getElementById('registerConfirmPassword').value;
        
        if (password !== confirmPassword) {
            alert('Пароли не совпадают!');
            return;
        }
        
        // Здесь должна быть логика регистрации
        // Для примера просто сохраняем данные пользователя
        currentUser = {
            name: name,
            email: email
        };
        
        registerModal.style.display = 'none';
        showUserAvatar();
    });
    
    // Показать аватар пользователя
    function showUserAvatar() {
        authButtons.style.display = 'none';
        userContainer.style.display = 'flex';
        userAvatar.textContent = currentUser.name.charAt(0).toUpperCase();
    }
    // Клик по профилю - показать/скрыть меню
    document.getElementById('userProfile').addEventListener('click', function(e) {
        e.stopPropagation();
        userMenu.style.display = userMenu.style.display === 'block' ? 'none' : 'block';
    });
    // Клик по кнопке "Профиль"
    profileBtn.addEventListener('click', function(e) {
        e.preventDefault();
        userMenu.style.display = 'none';
        alert(`Профиль пользователя:\nИмя: ${currentUser.name}\nEmail: ${currentUser.email}`);
    });
    
    // Клик по кнопке "Выход"
    logoutBtn.addEventListener('click', function(e) {
        e.preventDefault();
        currentUser = null;
        userContainer.style.display = 'none';
        authButtons.style.display = 'flex';
        userMenu.style.display = 'none';
    });
    
    // Закрытие меню при клике вне его
    document.addEventListener('click', function() {
        userMenu.style.display = 'none';
    });
});
const icon = document.querySelector('.icon');
const search = document.querySelector('.search');
const clear = document.querySelector('.clear');
icon.onclick = function()  {
search.classList.toggle('active');
};
clear.onclick = function()  {
    document.getElementById('mySearch').value='';
};