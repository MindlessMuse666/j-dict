(function() {
    console.log('[Mock Server] Инициализация мок-сервера...');

    const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

    const mockUser = {
        id: '1',
        email: 'demo@example.com',
        name: 'Demo User',
        role: 'user',
    };

    let mockWords = [
        {
            id: 1,
            jp: '猫',
            ru: 'кошка',
            kun: 'ねこ',
            on: 'ビョウ',
            ex_jp: '猫が好きです',
            ex_ru: 'Я люблю кошек',
            tags: ['animal', 'noun', 'n5'],
            kanji: '猫',
            reading: 'ねこ',
            romaji: 'neko',
            meanings: ['cat', 'кошка'],
            level: 'N5',
            created_at: new Date().toISOString()
        },
        {
            id: 2,
            jp: '犬',
            ru: 'собака',
            kun: 'いぬ',
            on: 'ケン',
            ex_jp: '犬が走っています',
            ex_ru: 'Собака бежит',
            tags: ['animal', 'noun', 'n5'],

            kanji: '犬',
            reading: 'いぬ',
            romaji: 'inu',
            meanings: ['dog', 'собака'],
            level: 'N5',
            created_at: new Date().toISOString()
        },
        {
            id: 3,
            jp: '食べる',
            ru: 'есть, кушать',
            kun: 'たべる',
            on: 'ショク',
            ex_jp: 'ご飯を食べる',
            ex_ru: 'Есть еду',
            tags: ['verb', 'n5'],

            kanji: '食べる',
            reading: 'たべる',
            romaji: 'taberu',
            meanings: ['to eat', 'есть'],
            level: 'N5',
            created_at: new Date().toISOString()
        },
        {
            id: 4,
            jp: '日本語',
            ru: 'японский язык',
            kun: 'にほんご',
            on: 'ニホンゴ',
            ex_jp: '日本語を勉強します',
            ex_ru: 'Я учу японский язык',
            tags: ['noun', 'language', 'n5'],

            kanji: '日本語',
            reading: 'にほんご',
            romaji: 'nihongo',
            meanings: ['Japanese language', 'японский язык'],
            level: 'N5',
            created_at: new Date().toISOString()
        },
        {
            id: 5,
            jp: '勉強',
            ru: 'учеба',
            kun: 'べんきょう',
            on: 'ベンキョウ',
            ex_jp: '毎日勉強します',
            ex_ru: 'Я учусь каждый день',
            tags: ['noun', 'suru-verb', 'n5'],

            kanji: '勉強',
            reading: 'べんきょう',
            romaji: 'benkyou',
            meanings: ['study', 'учеба'],
            level: 'N5',
            created_at: new Date().toISOString()
        }
    ];

    const originalXHR = window.XMLHttpRequest;

    function MockXHR() {
        const xhr = new originalXHR();
        const originalOpen = xhr.open;
        const originalSend = xhr.send;

        xhr.open = function(method, url, ...args) {
            this._method = method;
            this._url = url;
            return originalOpen.apply(this, [method, url, ...args]);
        };

        xhr.send = function(body) {
            const method = this._method.toUpperCase();
            const url = this._url;

            console.log(`[Mock Server] Intercepted: ${method} ${url}`);

            if (url.includes('/auth/login') && method === 'POST') {
                setTimeout(() => {
                    const response = {
                        data: {
                            user: mockUser,
                            token: 'mock-token-123'
                        }
                    };
                    Object.defineProperty(this, 'response', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'responseText', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'status', { value: 200 });
                    Object.defineProperty(this, 'readyState', { value: 4 });
                    this.dispatchEvent(new Event('readystatechange'));
                    this.dispatchEvent(new Event('load'));
                }, 500);
                return;
            }

            if (url.includes('/auth/me') && method === 'GET') {
                setTimeout(() => {
                    const response = {
                        data: {
                            user: mockUser
                        }
                    };
                    Object.defineProperty(this, 'response', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'responseText', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'status', { value: 200 });
                    Object.defineProperty(this, 'readyState', { value: 4 });
                    this.dispatchEvent(new Event('readystatechange'));
                    this.dispatchEvent(new Event('load'));
                }, 300);
                return;
            }

            if ((url.includes('/auth/me') || url.includes('/users/me') || url.includes('/users/avatar')) && (method === 'PATCH' || method === 'PUT')) {
                setTimeout(() => {
                    let parsedBody = {};
                    try {
                        parsedBody = JSON.parse(body);
                    } catch (e) {}

                    if (parsedBody.avatar_url) {
                        mockUser.avatar_url = parsedBody.avatar_url;
                    } else if (parsedBody.avatar) {
                        mockUser.avatar_url = parsedBody.avatar;
                    } else if (parsedBody.url) {
                        mockUser.avatar_url = parsedBody.url;
                    }

                    const response = {
                        data: {
                            user: mockUser,
                            success: true
                        }
                    };
                    Object.defineProperty(this, 'response', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'responseText', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'status', { value: 200 });
                    Object.defineProperty(this, 'readyState', { value: 4 });
                    this.dispatchEvent(new Event('readystatechange'));
                    this.dispatchEvent(new Event('load'));
                }, 500);
                return;
            }

            if ((url.includes('/users/avatar') || url.includes('/auth/avatar')) && method === 'POST') {
                 setTimeout(() => {
                    const response = {
                        data: {
                            user: mockUser,
                            success: true
                        }
                    };
                    Object.defineProperty(this, 'response', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'responseText', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'status', { value: 200 });
                    Object.defineProperty(this, 'readyState', { value: 4 });
                    this.dispatchEvent(new Event('readystatechange'));
                    this.dispatchEvent(new Event('load'));
                }, 800);
                return;
            }

            if (url.includes('/words') && method === 'GET') {
                setTimeout(() => {
                    const response = {
                        data: {
                            words: mockWords,
                            total_count: mockWords.length,
                            has_more: false,
                            next_cursor: 0
                        }
                    };
                    Object.defineProperty(this, 'response', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'responseText', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'status', { value: 200 });
                    Object.defineProperty(this, 'readyState', { value: 4 });
                    this.dispatchEvent(new Event('readystatechange'));
                    this.dispatchEvent(new Event('load'));
                }, 400);
                return;
            }
            
            if (url.includes('/words/search') && method === 'GET') {
                 setTimeout(() => {
                    const urlObj = new URL('http://dummy' + url);
                    const tags = urlObj.searchParams.get('tags');
                    const q = urlObj.searchParams.get('q') || '';
                    
                    let filtered = mockWords;
                    if (tags) {
                        const tagList = tags.split(',');
                        filtered = filtered.filter(w => w.tags.some(t => tagList.includes(t)));
                    }
                    
                    if (q) {
                         const qLower = q.toLowerCase();
                         filtered = filtered.filter(w => 
                             w.jp.includes(q) || 
                             w.kun.includes(q) || 
                             w.ru.toLowerCase().includes(qLower) ||
                             (w.meanings && w.meanings.some(m => m.toLowerCase().includes(qLower)))
                         );
                    }

                    const response = {
                        data: {
                            words: filtered,
                            total_count: filtered.length,
                            has_more: false,
                            next_cursor: 0
                        }
                    };
                    Object.defineProperty(this, 'response', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'responseText', { value: JSON.stringify(response) });
                    Object.defineProperty(this, 'status', { value: 200 });
                    Object.defineProperty(this, 'readyState', { value: 4 });
                    this.dispatchEvent(new Event('readystatechange'));
                    this.dispatchEvent(new Event('load'));
                }, 400);
                return;
            }

            return originalSend.apply(this, [body]);
        };

        return xhr;
    }

    window.XMLHttpRequest = MockXHR;
    console.log('[Mock Server] Initialized.');
})();
