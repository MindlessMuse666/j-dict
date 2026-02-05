const DEFAULT_USER = {
    name: 'Mr. Cat',
    email: 'cat@example.com',
    avatar: './assets/default_avatars/~mr~cat~.jpg',
    role: 'user'
};

const DEFAULT_WORDS = [
    {
        id: 1,
        japanese: '猫',
        reading: 'ねこ',
        russian: 'Кошка',
        tags: ['животные', 'милое']
    },
    {
        id: 2,
        japanese: '日本',
        reading: 'にほん',
        russian: 'Япония',
        tags: ['страна']
    },
    {
        id: 3,
        japanese: '勉強',
        reading: 'べんきょう',
        russian: 'Учёба',
        tags: ['действия']
    },
    {
        id: 4,
        japanese: '桜',
        reading: 'さくら',
        russian: 'Сакура',
        tags: ['природа', 'весна']
    },
    {
        id: 5,
        japanese: '美味しい',
        reading: 'おいしい',
        russian: 'Вкусный',
        tags: ['прилагательное', 'еда']
    }
];

const state = {
    user: JSON.parse(localStorage.getItem('user')) || null,
    words: JSON.parse(localStorage.getItem('words')) || DEFAULT_WORDS
};

const router = {
    routes: {
        '/': 'view-home',
        '/login': 'view-login',
        '/register': 'view-register',
        '/profile': 'view-profile'
    },

    init() {
        window.addEventListener('hashchange', () => this.handleRoute());
        this.handleRoute(); // Handle initial load
    },

    push(path) {
        window.location.hash = path;
    },

    handleRoute() {
        const hash = window.location.hash.slice(1) || '/';
        const viewId = this.routes[hash] || 'view-home';

        if (hash === '/' || hash === '/profile') {
            if (!state.user) {
                this.push('/login');
                return;
            }
        } else if (hash === '/login' || hash === '/register') {
            if (state.user) {
                this.push('/');
                return;
            }
        }

        document.querySelectorAll('main > div').forEach(el => el.classList.add('hidden'));
        
        const view = document.getElementById(viewId);
        if (view) {
            view.classList.remove('hidden');
            // Trigger specific view logic
            if (viewId === 'view-home') words.render();
            if (viewId === 'view-profile') profile.render();
        }

        auth.updateHeader();
    }
};

const auth = {
    login(email, password) {
        if (email === DEFAULT_USER.email) {
            this.setUser(DEFAULT_USER);
        } else {
            this.setUser({
                name: email.split('@')[0],
                email: email,
                avatar: './assets/default_avatars/~void~cat~.jpg',
                role: 'user'
            });
        }
        router.push('/');
    },

    register(name, email, password) {
        this.setUser({
            name: name,
            email: email,
            avatar: './assets/default_avatars/~lunishe~cat~.jpg',
            role: 'user'
        });
        router.push('/');
    },

    logout() {
        state.user = null;
        localStorage.removeItem('user');
        router.push('/login');
    },

    setUser(user) {
        state.user = user;
        localStorage.setItem('user', JSON.stringify(user));
    },

    updateHeader() {
        const navAuth = document.getElementById('nav-auth');
        const navGuest = document.getElementById('nav-guest');
        
        if (state.user) {
            navAuth.classList.remove('hidden');
            navGuest.classList.add('hidden');
            document.getElementById('user-name').textContent = state.user.name;
            document.getElementById('user-avatar').src = state.user.avatar;
        } else {
            navAuth.classList.add('hidden');
            navGuest.classList.remove('hidden');
        }
    }
};

const words = {
    render(filterText = '') {
        const grid = document.getElementById('words-grid');
        const emptyState = document.getElementById('empty-state');
        grid.innerHTML = '';

        const filteredWords = state.words.filter(word => {
            const search = filterText.toLowerCase();
            return word.japanese.includes(search) || 
                   word.reading.includes(search) || 
                   word.russian.toLowerCase().includes(search);
        });

        if (filteredWords.length === 0) {
            emptyState.classList.remove('hidden');
        } else {
            emptyState.classList.add('hidden');
            filteredWords.forEach(word => {
                const card = document.createElement('div');
                card.className = 'bg-white p-6 rounded-xl shadow-sm hover:shadow-md transition border border-gray-100 relative group';
                card.innerHTML = `
                    <div class="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition">
                        <button onclick="words.delete(${word.id})" class="text-gray-400 hover:text-red-500">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                        </button>
                    </div>
                    <div class="mb-2">
                        <span class="text-3xl font-bold text-gray-800 font-jp">${word.japanese}</span>
                        <span class="text-sm text-gray-500 ml-2 font-jp">【${word.reading}】</span>
                    </div>
                    <div class="text-lg text-gray-700 mb-4">${word.russian}</div>
                    <div class="flex flex-wrap gap-2">
                        ${word.tags.map(tag => `<span class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-md">#${tag}</span>`).join('')}
                    </div>
                `;
                grid.appendChild(card);
            });
        }
    },

    addWordPrompt() {
        const japanese = prompt('Японское слово:');
        if (!japanese) return;
        const reading = prompt('Чтение (хирагана/катакана):');
        const russian = prompt('Перевод:');
        
        if (japanese && reading && russian) {
            const newWord = {
                id: Date.now(),
                japanese,
                reading,
                russian,
                tags: ['новое']
            };
            state.words.unshift(newWord);
            localStorage.setItem('words', JSON.stringify(state.words));
            this.render(document.getElementById('search-input').value);
        }
    },

    delete(id) {
        if (confirm('Вы уверены, что хотите удалить это слово?')) {
            state.words = state.words.filter(w => w.id !== id);
            localStorage.setItem('words', JSON.stringify(state.words));
            this.render(document.getElementById('search-input').value);
        }
    },

    filter(text) {
        this.render(text);
    }
};

const profile = {
    render() {
        if (!state.user) return;
        document.getElementById('profile-name').textContent = state.user.name;
        document.getElementById('profile-email').textContent = state.user.email;
        document.getElementById('profile-avatar').src = state.user.avatar;
        document.getElementById('profile-word-count').textContent = state.words.length;
    }
};

document.addEventListener('DOMContentLoaded', () => {
    router.init();
});
