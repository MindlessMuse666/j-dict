/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                'onyomi': '#FEF3C7', // Amber-100 (warmer yellow)
                'kunyomi': '#E0F2FE', // Sky-100 (softer blue)
                'primary': '#BE123C', // Rose-700 (Traditional Red)
                'secondary': '#047857', // Emerald-700 (Bamboo Green)
                'danger': '#EF4444',
                'background': '#FAFAF9', // Stone-50 (Washi paper feel)
                'surface': '#FFFFFF',
                'text-main': '#1C1917', // Stone-900
                'text-muted': '#57534E', // Stone-600
            },
            fontFamily: {
                'sans': ['"Noto Sans"', '"Roboto"', 'sans-serif'],
                'jp': ['"Noto Sans JP"', '"Hiragino Sans"', '"Yu Gothic"', '"Meiryo"', 'sans-serif'],
            },
            typography: {
                jp: {
                    css: {
                        fontFamily: "'Noto Sans JP', 'Hiragino Sans', 'Yu Gothic', 'Meiryo', 'sans-serif'",
                        fontWeight: '400',
                    },
                },
            },
            backgroundImage: {
                'washi-pattern': "url(\"data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23d6d3d1' fill-opacity='0.1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E\")",
            }
        },
    },
    plugins: [],
}