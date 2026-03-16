/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#3157d5',
          secondary: '#0f8b8d',
          surface: '#ffffff',
          background: '#f4f7fb',
          success: '#1ea971',
          warning: '#d68b1f',
          error: '#d24c59',
        },
      },
    },
  },
})
