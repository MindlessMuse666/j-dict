(function() {
    console.log('[Mock Server] Initializing mock server...');

    const mockUser = {
        id: '1',
        email: 'demo@example.com',
        name: 'Demo User',
        role: 'user',
        avatar_url: '/assets/default_avatars/~mr~cat~.jpg'
    };

    let mockWords = [
        {
            id: 1,
            kanji: '猫',
            reading: 'ねこ',
            romaji: 'neko',
            meanings: ['cat'],
            level: 'N5',
            tags: ['animal', 'noun'],
            created_at: new Date().toISOString()
        },
        {
            id: 2,
            kanji: '犬',
            reading: 'いぬ',
            romaji: 'inu',
            meanings: ['dog'],
            level: 'N5',
            tags: ['animal', 'noun'],
            created_at: new Date().toISOString()
        },
        {
            id: 3,
            kanji: '食べる',
            reading: 'たべる',
            romaji: 'taberu',
            meanings: ['to eat'],
            level: 'N5',
            tags: ['verb'],
            created_at: new Date().toISOString()
        },
        {
            id: 4,
            kanji: '日本語',
            reading: 'にほんご',
            romaji: 'nihongo',
            meanings: ['Japanese language'],
            level: 'N5',
            tags: ['noun'],
            created_at: new Date().toISOString()
        },
        {
            id: 5,
            kanji: '勉強',
            reading: 'べんきょう',
            romaji: 'benkyou',
            meanings: ['study'],
            level: 'N5',
            tags: ['noun', 'suru-verb'],
            created_at: new Date().toISOString()
        }
    ];

    // Hijack XMLHttpRequest
    const originalXHR = window.XMLHttpRequest;

    function MockXHR() {
        const xhr = new originalXHR();
        this._xhr = xhr;
        this._headers = {};
        this.responseType = '';
        this.response = null;
        this.responseText = null;
        this.status = 0;
        this.statusText = '';
        this.readyState = 0;
        this.onreadystatechange = null;
        this.onload = null;
        this.onerror = null;

        const self = this;

        // Proxy common properties/methods
        Object.defineProperty(this, 'upload', {
            get: () => xhr.upload
        });
    }

    MockXHR.prototype.open = function(method, url, async, user, password) {
        this._method = method.toUpperCase();
        this._url = url;
        // Check if URL is internal API call (starts with /api or is relative)
        // Adjust logic to match axios base URL
        this._isApi = url.includes('/api') || (!url.startsWith('http') && !url.includes('.js') && !url.includes('.css'));
        
        if (!this._isApi) {
            return this._xhr.open(method, url, async, user, password);
        }
    };

    MockXHR.prototype.setRequestHeader = function(header, value) {
        this._headers[header] = value;
        if (!this._isApi) {
            this._xhr.setRequestHeader(header, value);
        }
    };

    MockXHR.prototype.send = function(data) {
        if (!this._isApi) {
            const self = this;
            this._xhr.onreadystatechange = function() {
                self.readyState = this.readyState;
                self.status = this.status;
                self.statusText = this.statusText;
                self.response = this.response;
                self.responseText = this.responseText;
                if (self.onreadystatechange) self.onreadystatechange();
                if (this.readyState === 4 && self.onload) self.onload();
            };
            this._xhr.onerror = function() {
                if (self.onerror) self.onerror();
            };
            return this._xhr.send(data);
        }

        console.log(`[Mock Server] ${this._method} ${this._url}`, data);

        // Simulate network delay
        setTimeout(() => {
            let response = null;
            let status = 404;

            try {
                // Normalize URL to remove base /api if present
                const path = this._url.replace('/api', '');

                // Routes
                if (path === '/auth/login' && this._method === 'POST') {
                    const { email, password } = JSON.parse(data);
                    if (email === 'demo@example.com' && password === 'password') {
                        response = {
                            data: {
                                user: mockUser,
                                token: 'mock-token-123'
                            }
                        };
                        status = 200;
                    } else {
                        response = { error: 'Invalid credentials' };
                        status = 401;
                    }
                } else if (path === '/auth/register' && this._method === 'POST') {
                    response = {
                        data: {
                            user: mockUser,
                            token: 'mock-token-123'
                        }
                    };
                    status = 200;
                } else if (path === '/auth/me' && this._method === 'GET') {
                    if (this._headers['Authorization'] === 'Bearer mock-token-123') {
                        response = { data: mockUser };
                        status = 200;
                    } else {
                        response = { error: 'Unauthorized' };
                        status = 401;
                    }
                } else if (path === '/auth/logout' && this._method === 'POST') {
                    response = { success: true };
                    status = 200;
                } else if (path.startsWith('/words') && this._method === 'GET') {
                    if (path.includes('/search')) {
                         // Parse query params basic implementation
                        const urlObj = new URL('http://localhost' + path); // Dummy host for parsing
                        const q = urlObj.searchParams.get('query')?.toLowerCase() || '';
                        
                        let results = mockWords;
                        if (q) {
                            results = results.filter(w => 
                                w.kanji.includes(q) || 
                                w.reading.includes(q) || 
                                w.meanings.some(m => m.toLowerCase().includes(q))
                            );
                        }
                         response = {
                            data: {
                                words: results,
                                next_cursor: 0,
                                has_more: false,
                                total_count: results.length
                            }
                        };
                        status = 200;
                    } else {
                        response = {
                            data: {
                                words: mockWords,
                                next_cursor: 0,
                                has_more: false,
                                total_count: mockWords.length
                            }
                        };
                        status = 200;
                    }
                } else if (path === '/words' && this._method === 'POST') {
                    const newWord = JSON.parse(data);
                    const created = {
                        ...newWord,
                        id: mockWords.length + 1,
                        created_at: new Date().toISOString()
                    };
                    mockWords.unshift(created);
                    response = { data: created };
                    status = 201;
                } else if (path.match(/\/words\/\d+/) && this._method === 'DELETE') {
                     const id = parseInt(path.split('/').pop());
                     mockWords = mockWords.filter(w => w.id !== id);
                     response = { success: true };
                     status = 200;
                } else if (path.match(/\/words\/\d+/) && this._method === 'PUT') {
                     const id = parseInt(path.split('/').pop());
                     const updateData = JSON.parse(data);
                     const index = mockWords.findIndex(w => w.id === id);
                     if (index !== -1) {
                        mockWords[index] = { ...mockWords[index], ...updateData };
                        response = { data: mockWords[index] };
                        status = 200;
                     } else {
                         status = 404;
                     }
                }
            } catch (e) {
                console.error('[Mock Server] Error:', e);
                status = 500;
                response = { error: e.message };
            }

            // Respond
            this.readyState = 4;
            this.status = status;
            this.statusText = status === 200 ? 'OK' : 'Error';
            this.response = JSON.stringify(response); // Axios expects string if responseType is text/json
            this.responseText = this.response;

            if (this.onreadystatechange) this.onreadystatechange();
            if (this.onload) this.onload();

        }, 300);
    };
    
    // Add missing methods that might be called
    MockXHR.prototype.abort = function() {};
    MockXHR.prototype.getAllResponseHeaders = function() { return ''; };
    MockXHR.prototype.getResponseHeader = function(h) { return null; };

    // Overwrite global
    window.XMLHttpRequest = MockXHR;
})();
