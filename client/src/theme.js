// theme.js
import { extendTheme } from '@chakra-ui/react';

const theme = extendTheme({
  fonts: {
    heading: `'Merriweather', serif`,
    body: `'Merriweather', serif`,
  },
  config: {
    initialColorMode: 'dark',
    useSystemColorMode: false,
  },
  colors: {
    // You can customize these colors to your preference
    bg: {
      primary: '#2f3349', // Dark background
      secondary: '#2D3748', // Slightly lighter background for contrast
    },
    text: {
      primary: '#E2E8F0', // Light text for dark background
      secondary: '#A0AEC0', // Slightly darker text for less emphasis
    },
    accent: {
      primary: '#4FD1C5', // Teal accent color
      secondary: '#38B2AC', // Darker teal for hover states
    },
  },
  styles: {
    global: {
      body: {
        bg: 'bg.primary',
        color: 'text.primary',
      },
    },
  },
});

export default theme;