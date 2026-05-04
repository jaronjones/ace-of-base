# DESIGN.md - NEON_OS

## Brand Identity
NEON_OS is a high-fidelity retro-future dashboard system. It blends 1980s synthwave aesthetics (palm trees, neon sunsets, CRT textures) with modern high-density data visualization.

### Visual Principles
- **Atmospheric Contrast**: Deep black backgrounds (#000000 or near-black) paired with high-luminance neon accents.
- **Digital Nostalgia**: Use of horizontal scan-lines, glowing interactive elements, and sharp, monospaced-style typography.
- **Cybernetic Hierarchy**: Information is organized into modular "terminal" blocks with subtle translucent backdrops and glowing borders.

---

## Color Systems

The project supports several distinct color palettes derived from the brand's evolution.

### 1. Synth-Grid (The Original)
- **Primary**: #FF00FF (Electric Magenta)
- **Secondary**: #00FFFF (Cyber Cyan)
- **Surface**: #131313 (Matte Black)
- **Accents**: Violet-600, Sky-400
*The foundation of the NEON_OS experience, balancing classic vaporwave tones with a sharp, high-contrast grid aesthetic.*

### 2. Neon Horizon (Core Aesthetic)
- **Primary**: #FF2DED (Magenta Glow)
- **Secondary**: #22D3EE (Cyan Data)
- **Surface**: #1C0F19 (Dark Plum/Black)
- **Accents**: Indigo-500, Purple-500

### 3. Solaris Terminal
- **Primary**: #FF4D2D (Solar Red)
- **Secondary**: #FFCC00 (Amber Warning)
- **Surface**: #1F0F0C (Deep Umber)

### 4. Retro-Future
- **Primary**: #A855F7 (Purple)
- **Secondary**: #F97316 (Orange Sunset)
- **Surface**: #131313 (Matte Black)

### 5. Tech-Noir
- **Primary**: #E63946 (Crimson)
- **Secondary**: #10B981 (Emerald/Teal)
- **Surface**: #131313

### 6. Cyan-Sunset
- **Primary**: #00F5FF (Vibrant Cyan)
- **Secondary**: #F59E0B (Sunset Amber)
- **Surface**: #131313

### 7. Sky-Blue
- **Primary**: #7DD3FC (Sky Blue)
- **Secondary**: #C084FC (Violet)
- **Surface**: #131313

---

## Typography

### Primary Font: Space Grotesk
- **Usage**: Headlines, Navigation, Brand Marks.
- **Style**: Bold, uppercase, tracking-widest for that "computer terminal" feel.

### Secondary Font: Monospace (Standard System)
- **Usage**: System logs, technical data, widget labels.
- **Style**: High-readability, fixed-width.

---

## Component Guidelines

### Navigation
- **Desktop**: Fixed left side-nav with active-state glow and border-left highlights.
- **Mobile**: Bottom navigation bar with large icon targets and neon-pulsing active icons.

### Widgets & Containers
- **Shapes**: Roundness set to `ROUND_FOUR` (4px) or `UNSPECIFIED` (Sharp corners).
- **Borders**: 1px or 2px solid borders using primary/secondary neon colors at 30-50% opacity.
- **Shadows**: Intense outer glow (`drop-shadow`) on primary interactive elements.

---

## Assets
- **Logo**: Criss-crossing palm trees silhouetted against a horizontal-lined setting sun.
- **Icons**: Outlined Material Symbols, styled with neon color weight.
