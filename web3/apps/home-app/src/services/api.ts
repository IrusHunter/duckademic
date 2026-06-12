import axios from 'axios'

// Єдиний axios-інстанс на весь home-app.
// Токен береться з localStorage — туди його кладе shell після логіну
// (shell — єдиний власник токенів, home-app їх лише читає).
//
// Коли модуль працює всередині shell, активний і глобальний інтерсептор shell,
// але власний інтерсептор гарантує підстановку токена навіть при ізольованому
// запуску модуля. Відносні URL ('/api/...') резолвляться відносно origin shell,
// де налаштований проксі на api-gateway.
export const api = axios.create({ withCredentials: true })

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
