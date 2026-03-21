package web

// officeHTML contains the 8-bit pixel art office dashboard
var officeHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Agent House - 8-Bit Office</title>
    <link href="https://fonts.googleapis.com/css2?family=Press+Start+2P&display=swap" rel="stylesheet">
    <style>
        :root {
            --color-ceo: #4A90A4;
            --color-pm: #7B68EE;
            --color-ux: #FF6B6B;
            --color-ui: #4ECDC4;
            --color-security: #F39C12;
            --color-architect: #9B59B6;
            --color-senior-dev: #2ECC71;
            --color-junior-dev: #3498DB;

            --bg-app: #1a1025;
            --bg-panel: #1e1430;
            --bg-dark: #140e1e;
            --text-primary: #e8e0f0;
            --text-secondary: #9088a0;
            --text-dim: #5a5268;
            --border-color: #3a2850;
            --accent: #4ECDC4;

            --floor-main: #4a3d30;
            --floor-alt: #6b5a42;
            --wall-main: #2a2040;
            --wall-light: #342850;
        }

        *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'Press Start 2P', monospace;
            background: var(--bg-app);
            color: var(--text-primary);
            min-height: 100vh;
            overflow: hidden;
            image-rendering: pixelated;
        }

        @media (prefers-reduced-motion: reduce) {
            *, *::before, *::after {
                animation-duration: 0.01ms !important;
                transition-duration: 0.01ms !important;
            }
        }

        .app { display: flex; flex-direction: column; height: 100vh; }

        /* ======================== HEADER ======================== */
        .header {
            display: flex; align-items: center; justify-content: space-between;
            padding: 0 16px; background: var(--bg-dark); border-bottom: 4px solid var(--border-color);
            height: 44px; z-index: 100; flex-shrink: 0;
        }
        .header-left { display: flex; align-items: center; gap: 16px; }
        .logo { font-size: 10px; color: var(--accent); letter-spacing: 2px; }
        .phase-pill {
            display: flex; align-items: center; gap: 6px; padding: 3px 10px;
            background: var(--bg-panel); border: 2px solid var(--border-color); font-size: 7px;
            color: var(--text-secondary); letter-spacing: 1px;
        }
        .phase-dot { width: 6px; height: 6px; background: var(--accent); animation: blink 1s step-end infinite; }
        @keyframes blink { 50% { opacity: 0; } }
        .header-right { display: flex; align-items: center; gap: 10px; }
        .view-toggle {
            padding: 3px 10px; background: transparent; border: 2px solid var(--border-color);
            color: var(--text-secondary); font-family: inherit; font-size: 7px; cursor: pointer;
            letter-spacing: 1px; text-decoration: none;
        }
        .view-toggle:hover { border-color: var(--accent); color: var(--accent); }
        .ws-dot { width: 8px; height: 8px; background: #555; }
        .ws-dot.on { background: #2ECC71; }
        .project-sel {
            background: var(--bg-dark); border: 2px solid var(--border-color); color: var(--text-primary);
            font-family: inherit; font-size: 7px; padding: 3px 6px; cursor: pointer;
        }

        .new-proj-btn {
            padding: 3px 8px; background: transparent; border: 2px solid var(--border-color);
            color: var(--accent); font-family: inherit; font-size: 10px; cursor: pointer;
            line-height: 1;
        }
        .new-proj-btn:hover { background: var(--accent); color: var(--bg-dark); }

        /* New project modal */
        .np-modal {
            position: fixed; top: 0; left: 0; right: 0; bottom: 0;
            background: rgba(10,6,14,0.85); z-index: 200;
            display: flex; align-items: center; justify-content: center;
            animation: fadeIn 0.2s;
        }
        .np-box {
            background: var(--bg-panel); border: 4px solid var(--accent);
            padding: 16px; width: 320px; image-rendering: auto;
        }
        .np-title {
            font-size: 10px; color: var(--accent); letter-spacing: 2px;
            margin-bottom: 12px; text-align: center;
        }
        .np-input {
            width: 100%; background: var(--bg-dark); border: 2px solid var(--border-color);
            color: var(--text-primary); font-family: inherit; font-size: 9px;
            padding: 8px 10px; outline: none; margin-bottom: 10px;
        }
        .np-input:focus { border-color: var(--accent); }
        .np-input::placeholder { color: var(--text-dim); }
        .np-hint {
            font-size: 6px; color: var(--text-dim); margin-bottom: 12px;
            letter-spacing: 1px;
        }
        .np-actions { display: flex; gap: 8px; justify-content: flex-end; }
        .np-btn {
            padding: 6px 14px; font-family: inherit; font-size: 8px;
            cursor: pointer; letter-spacing: 1px; border: 2px solid;
        }
        .np-btn.cancel {
            background: transparent; border-color: var(--border-color);
            color: var(--text-secondary);
        }
        .np-btn.cancel:hover { border-color: #E74C3C; color: #E74C3C; }
        .np-btn.create {
            background: var(--accent); border-color: #3aaa9a;
            color: var(--bg-dark);
        }
        .np-btn.create:hover { background: #5de0d6; }
        .np-error { font-size: 7px; color: #E74C3C; margin-bottom: 8px; display: none; }

        /* ======================== MAIN AREA ======================== */
        .main-area { display: flex; flex: 1; overflow: hidden; }

        /* ======================== OFFICE ======================== */
        .office-wrap {
            flex: 1; position: relative; overflow: hidden;
            background: var(--wall-main);
        }
        .office-scene {
            position: absolute; top: 50%; left: 50%;
            transform: translate(-50%, -50%);
            width: 1000px; height: 520px;
        }

        /* ---- Wall ---- */
        .wall {
            position: absolute; top: 0; left: 0; right: 0; height: 130px;
            background: linear-gradient(180deg, #221838 0%, var(--wall-main) 100%);
            border-bottom: 4px solid #443060;
        }

        /* Windows */
        .window {
            position: absolute; top: 10px; width: 56px; height: 56px;
            background: #0a1628; border: 4px solid #554070;
            box-shadow: inset 0 0 20px rgba(100,160,255,0.1);
        }
        .window.w1 { left: 40px; }
        .window.w2 { left: 160px; }
        .window.w3 { right: 160px; }
        .window.w4 { right: 40px; }
        .window-pane {
            position: absolute; top: 0; left: 50%; width: 2px; height: 100%; background: #554070;
        }
        .window-pane.h {
            top: 50%; left: 0; width: 100%; height: 2px;
        }
        .star {
            position: absolute; width: 2px; height: 2px; background: #fff;
            animation: twinkle 3s ease-in-out infinite;
        }
        @keyframes twinkle { 0%,100%{opacity:1}50%{opacity:0.2} }

        /* Wall Posters */
        .poster {
            position: absolute; border: 3px solid #554070; padding: 4px;
            font-size: 8px; text-align: center; line-height: 1.4;
            background: var(--wall-light);
        }
        .poster .poster-art {
            display: block; margin-bottom: 2px; font-size: 14px; line-height: 1;
        }
        .poster .poster-text {
            display: block; color: #bba8d0; letter-spacing: 1px; text-transform: uppercase;
        }
        .poster.p1 { top: 16px; left: 260px; width: 68px; }
        .poster.p2 { top: 14px; left: 350px; width: 64px; }
        .poster.p3 { top: 16px; right: 260px; width: 64px; }

        /* Whiteboard */
        .whiteboard {
            position: absolute; top: 6px; left: 50%; transform: translateX(-50%);
            width: 150px; height: 86px; background: #eee8dd; border: 4px solid #888;
            cursor: pointer; z-index: 10; transition: box-shadow 0.2s;
        }
        .whiteboard:hover { box-shadow: 0 0 16px rgba(78,205,196,0.3); }
        .wb-header {
            background: #cc4444; color: #fff; font-size: 7px; text-align: center;
            padding: 2px; letter-spacing: 1px;
        }
        .wb-body {
            padding: 4px; font-size: 7px; color: #444; line-height: 1.5;
            font-family: 'Courier New', monospace;
        }
        .wb-line { display: block; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
        .wb-line::before { content: '> '; color: #999; }

        /* ---- Floor ---- */
        .floor {
            position: absolute; top: 130px; left: 0; right: 0; bottom: 0;
            background:
                repeating-conic-gradient(var(--floor-main) 0% 25%, var(--floor-alt) 0% 50%) 0 0 / 48px 48px;
        }

        /* Zone Carpets */
        .zone {
            position: absolute; border-radius: 6px; z-index: 1;
        }
        .zone-label {
            position: absolute; top: -14px; left: 50%; transform: translateX(-50%);
            font-size: 8px; letter-spacing: 2px; padding: 2px 8px;
            white-space: nowrap; border: 2px solid; text-transform: uppercase;
        }
        .zone.exec {
            left: 20px; top: 16px; width: 250px; height: 200px;
            background: rgba(74,144,164,0.14); border: 2px solid rgba(74,144,164,0.25);
        }
        .zone.exec .zone-label { color: var(--color-ceo); border-color: rgba(74,144,164,0.4); background: rgba(74,144,164,0.15); }
        .zone.design {
            left: 290px; top: 16px; width: 300px; height: 200px;
            background: rgba(255,107,107,0.12); border: 2px solid rgba(255,107,107,0.2);
        }
        .zone.design .zone-label { color: var(--color-ux); border-color: rgba(255,107,107,0.35); background: rgba(255,107,107,0.12); }
        .zone.tech {
            left: 610px; top: 16px; width: 370px; height: 200px;
            background: rgba(155,89,182,0.12); border: 2px solid rgba(155,89,182,0.2);
        }
        .zone.tech .zone-label { color: var(--color-architect); border-color: rgba(155,89,182,0.35); background: rgba(155,89,182,0.12); }
        .zone.dev {
            left: 50%; transform: translateX(-50%); top: 170px; width: 380px; height: 150px;
            background: rgba(46,204,113,0.12); border: 2px solid rgba(46,204,113,0.2);
        }
        .zone.dev .zone-label { color: var(--color-senior-dev); border-color: rgba(46,204,113,0.35); background: rgba(46,204,113,0.12); }

        /* ---- Furniture ---- */
        .desk {
            position: absolute; z-index: 5;
        }
        .desk-top {
            width: 100px; height: 44px;
            background: linear-gradient(180deg, #a08050 0%, #8a6a3c 100%);
            border: 3px solid #705828;
            border-top: 3px solid #b89860;
            position: relative;
        }
        .desk-top.large { width: 120px; }
        .desk-leg {
            position: absolute; bottom: -8px; width: 6px; height: 8px; background: #705828;
        }
        .desk-leg.dl { left: 6px; }
        .desk-leg.dr { right: 6px; }

        /* (nameplate moved to character front) */

        /* Monitor (back view - user sees rear of monitor) */
        .monitor-back {
            position: absolute; top: -20px; left: 50%; transform: translateX(-50%);
            width: 44px; height: 20px; background: #2a2a3a; border: 3px solid #3a3a4a;
        }
        .monitor-back-panel {
            width: 100%; height: 100%; background: #222233;
        }
        /* Screen glow spills over the top when active */
        .monitor-back.on {
            box-shadow: 0 -6px 16px rgba(51,255,51,0.15);
        }
        .monitor-stand {
            width: 12px; height: 4px; background: #3a3a4a; margin: 0 auto;
        }

        /* Nameplate with status dot (below desk, facing user) */
        .nameplate-front {
            position: absolute; bottom: -22px; left: 50%; transform: translateX(-50%);
            padding: 2px 8px 2px 6px; font-size: 8px; letter-spacing: 2px;
            border: 2px solid; white-space: nowrap; z-index: 22;
            display: flex; align-items: center; gap: 5px;
        }
        .status-dot {
            width: 6px; height: 6px; flex-shrink: 0;
        }
        .status-dot.idle { background: #555; }
        .status-dot.working {
            background: #2ECC71;
            animation: dotPulseGlow 1.2s ease-in-out infinite;
        }
        .status-dot.error {
            background: #E74C3C;
            animation: dotPulseGlow 0.4s step-end infinite;
        }
        .status-dot.completed {
            background: var(--accent);
        }
        @keyframes dotPulseGlow {
            0%,100%{opacity:1;box-shadow:0 0 4px currentColor}
            50%{opacity:0.4;box-shadow:none}
        }

        /* Desk items */
        .desk-item {
            position: absolute; z-index: 7;
        }
        .rubber-duck {
            width: 18px; height: 14px; background: #FFE000; border-radius: 50% 50% 40% 40%;
            border: 2px solid #DAA520; position: relative;
        }
        .rubber-duck::before {
            content: ''; position: absolute; top: 3px; left: 12px;
            width: 8px; height: 5px; background: #FFA500; border-radius: 0 50% 50% 0;
        }
        .rubber-duck::after {
            content: ''; position: absolute; top: 2px; left: 4px;
            width: 3px; height: 3px; background: #333; border-radius: 50%;
        }
        .coffee-cup {
            width: 14px; height: 14px; background: #fff; border: 2px solid #ddd;
            border-radius: 0 0 3px 3px; position: relative;
        }
        .coffee-cup::before {
            content: ''; position: absolute; top: -5px; left: 2px;
            width: 8px; height: 5px; background: #6B3410; border-radius: 3px 3px 0 0;
        }
        .coffee-cup::after {
            content: '~'; position: absolute; top: -14px; left: 2px;
            font-size: 10px; color: rgba(255,255,255,0.4);
            animation: steam 2s ease-in-out infinite;
        }
        @keyframes steam {
            0%,100%{opacity:0.2;transform:translateY(0)}
            50%{opacity:0.6;transform:translateY(-3px)}
        }

        /* Sleeping cat */
        .pixel-cat {
            position: absolute; z-index: 8; cursor: pointer;
        }
        .cat-body {
            width: 30px; height: 14px; background: #E87040; border-radius: 10px 10px 4px 4px;
            position: relative;
        }
        .cat-body::before {
            content: ''; position: absolute; top: -5px; left: 3px;
            border-left: 5px solid transparent; border-right: 5px solid transparent;
            border-bottom: 6px solid #E87040;
        }
        .cat-body::after {
            content: ''; position: absolute; top: -5px; right: 3px;
            border-left: 5px solid transparent; border-right: 5px solid transparent;
            border-bottom: 6px solid #E87040;
        }
        .cat-tail {
            position: absolute; bottom: 3px; right: -12px;
            width: 14px; height: 4px; background: #E87040; border-radius: 0 4px 4px 0;
            animation: tailWag 3s ease-in-out infinite;
        }
        @keyframes tailWag {
            0%,100%{transform:rotate(-5deg)}50%{transform:rotate(10deg)}
        }
        .cat-zzz {
            position: absolute; top: -14px; right: -4px; font-size: 8px;
            color: var(--text-dim); animation: catSleep 4s ease-in-out infinite;
        }
        @keyframes catSleep {
            0%,100%{opacity:0;transform:translateY(0)}
            50%{opacity:0.8;transform:translateY(-6px)}
        }

        /* Coffee Machine */
        .coffee-machine {
            position: absolute; z-index: 5; cursor: pointer;
        }
        .cm-body {
            width: 28px; height: 36px; background: #666; border: 3px solid #555;
            position: relative;
        }
        .cm-top {
            width: 32px; height: 6px; background: #777; border: 2px solid #555;
            margin-left: -2px;
        }
        .cm-screen {
            width: 14px; height: 8px; background: #0a3a0a; margin: 4px auto 0;
            border: 1px solid #444;
        }
        .cm-spout {
            width: 4px; height: 6px; background: #555; margin: 2px auto 0;
        }
        .cm-cup {
            width: 10px; height: 8px; background: #fff; border: 1px solid #ddd;
            margin: 0 auto; border-radius: 0 0 2px 2px;
        }
        .cm-label {
            font-size: 4px; text-align: center; color: var(--text-dim);
            margin-top: 4px; letter-spacing: 1px;
        }

        /* Bookshelf */
        .bookshelf {
            position: absolute; z-index: 4;
        }
        .shelf-frame {
            width: 40px; height: 52px; background: #705828; border: 3px solid #5a4420;
            display: flex; flex-direction: column; padding: 3px; gap: 2px;
        }
        .shelf-row {
            display: flex; gap: 1px; flex: 1; align-items: flex-end;
        }
        .book {
            width: 5px; border-radius: 1px 1px 0 0;
        }

        /* Plants */
        .plant {
            position: absolute; z-index: 4;
        }
        .plant-pot {
            width: 18px; height: 12px; background: #c2703a;
            border: 2px solid #a05a2a; border-radius: 0 0 4px 4px;
        }
        .plant-top {
            width: 24px; height: 18px; margin: -6px -3px 0;
            background: radial-gradient(ellipse, #32a852 60%, #228B22 100%);
            border-radius: 50%;
        }
        .plant-top::before {
            content: ''; position: absolute; top: -4px; left: 6px;
            width: 12px; height: 10px; background: #3bc45f; border-radius: 50%;
        }

        /* ======================== CHARACTERS ======================== */
        .agent {
            position: absolute; z-index: 20; cursor: pointer;
            transition: transform 0.15s;
        }
        .agent:hover { transform: scale(1.1); z-index: 25; }
        .agent:hover .agent-label { opacity: 1; }

        .char {
            width: 36px; height: 52px; position: relative;
            image-rendering: pixelated;
        }

        /* Hair (top of head) */
        .hair {
            position: absolute; z-index: 3;
        }

        /* Head */
        .head {
            width: 22px; height: 20px; position: absolute; top: 6px; left: 7px;
            z-index: 2; border: 2px solid rgba(0,0,0,0.15);
        }

        /* Face */
        .face { position: absolute; top: 8px; left: 4px; width: 14px; height: 8px; z-index: 3; }
        .eye {
            position: absolute; top: 0; width: 3px; height: 3px; background: #2a1a0a;
        }
        .eye.l { left: 1px; }
        .eye.r { right: 1px; }
        .mouth {
            position: absolute; bottom: 0; left: 50%; transform: translateX(-50%);
            width: 4px; height: 2px; background: #c47a5a; border-radius: 0 0 2px 2px;
        }

        /* Blink */
        .blink .eye { height: 1px; margin-top: 1px; border-radius: 0; }

        /* Body */
        .torso {
            position: absolute; top: 26px; left: 6px; width: 24px; height: 14px;
            z-index: 1; border: 2px solid rgba(0,0,0,0.15);
        }

        /* Arms */
        .arm {
            position: absolute; top: 28px; width: 8px; height: 12px;
            z-index: 1; border: 2px solid rgba(0,0,0,0.12);
        }
        .arm.la { left: 0; border-radius: 4px 0 0 4px; }
        .arm.ra { right: 0; border-radius: 0 4px 4px 0; }

        /* Typing */
        .typing .arm.la { animation: typeArmL 0.3s step-end infinite; }
        .typing .arm.ra { animation: typeArmR 0.3s step-end infinite 0.15s; }
        @keyframes typeArmL { 50%{transform:translateY(-3px) rotate(-5deg)} }
        @keyframes typeArmR { 50%{transform:translateY(-3px) rotate(5deg)} }

        /* Legs */
        .legs {
            position: absolute; top: 40px; left: 9px; display: flex; gap: 2px; z-index: 0;
        }
        .leg { width: 7px; height: 12px; background: #3a3552; border: 1px solid #2a2542; }

        /* ---- Hair Styles ---- */
        /* CEO: slicked gray executive */
        .hair-ceo {
            top: 0; left: 5px; width: 26px; height: 12px;
            background: #8a8a90; border-radius: 6px 6px 0 0;
            border: 2px solid #6a6a70; border-bottom: none;
        }
        /* PM: neat brown side-part */
        .hair-pm {
            top: 0; left: 6px; width: 24px; height: 14px;
            background: #6B3A20; border-radius: 8px 4px 0 0;
            border: 2px solid #4a2810; border-bottom: none;
        }
        /* UX: pink bob */
        .hair-ux {
            top: 0; left: 4px; width: 28px; height: 16px;
            background: #e06090; border-radius: 10px 10px 4px 4px;
            border: 2px solid #c04878; border-bottom: none;
        }
        /* UI: blue messy/creative */
        .hair-ui {
            top: -2px; left: 4px; width: 28px; height: 16px;
            background: #40b0c0; border-radius: 6px 12px 2px 6px;
            border: 2px solid #2a8a9a; border-bottom: none;
        }
        /* Security: buzz cut dark */
        .hair-sec {
            top: 2px; left: 7px; width: 22px; height: 8px;
            background: #3a3a40; border-radius: 4px 4px 0 0;
            border: 2px solid #2a2a30; border-bottom: none;
        }
        /* Architect: neat with glasses */
        .hair-arch {
            top: 0; left: 6px; width: 24px; height: 12px;
            background: #4a3520; border-radius: 6px 6px 0 0;
            border: 2px solid #3a2510; border-bottom: none;
        }
        .glasses {
            position: absolute; top: 8px; left: 3px; z-index: 4;
            width: 16px; height: 6px;
            border: 2px solid #8888aa; border-radius: 3px;
        }
        .glasses::before {
            content: ''; position: absolute; top: 0; left: 50%;
            transform: translateX(-50%); width: 3px; height: 2px;
            border-top: 2px solid #8888aa;
        }
        /* Senior Dev: ginger beard */
        .hair-sr {
            top: 0; left: 5px; width: 26px; height: 14px;
            background: #c45020; border-radius: 8px 8px 0 0;
            border: 2px solid #a03a10; border-bottom: none;
        }
        .beard {
            position: absolute; top: 18px; left: 5px; z-index: 4;
            width: 16px; height: 8px; background: #c45020;
            border-radius: 0 0 6px 6px; border: 1px solid #a03a10;
            margin-left: 3px;
        }
        /* Junior Dev: cap */
        .hair-jr {
            top: -2px; left: 4px; width: 28px; height: 10px;
            background: var(--color-junior-dev);
            border-radius: 4px 4px 0 0;
            border: 2px solid #2a7ab0;
        }
        .cap-brim {
            position: absolute; top: 4px; left: 1px; z-index: 4;
            width: 16px; height: 4px; background: var(--color-junior-dev);
            border-radius: 2px; border: 1px solid #2a7ab0;
        }

        /* ---- Role Accessories ---- */
        .acc {
            position: absolute; z-index: 5; image-rendering: pixelated;
        }
        /* CEO: gold badge/tie */
        .acc-ceo {
            top: 28px; left: 14px; width: 8px; height: 6px;
            background: #FFD700; clip-path: polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%);
        }
        /* PM: clipboard */
        .acc-pm {
            top: 26px; right: -4px; width: 10px; height: 14px;
            background: #e8dcc8; border: 1px solid #baa888;
        }
        .acc-pm::before {
            content: ''; position: absolute; top: -2px; left: 2px;
            width: 6px; height: 3px; background: #888; border-radius: 1px;
        }
        .acc-pm::after {
            content: ''; position: absolute; top: 4px; left: 2px;
            width: 6px; height: 1px; background: #aaa;
            box-shadow: 0 3px 0 #aaa, 0 6px 0 #aaa;
        }
        /* UX: heart icon (empathy) */
        .acc-ux {
            top: 0px; right: -6px; width: 10px; height: 9px;
            background: var(--color-ux);
            clip-path: polygon(50% 100%, 0% 35%, 10% 0%, 40% 0%, 50% 20%, 60% 0%, 90% 0%, 100% 35%);
        }
        /* UI: palette */
        .acc-ui {
            top: 26px; left: -8px; width: 14px; height: 12px;
            background: #ddd; border-radius: 50% 50% 50% 20%;
            border: 1px solid #bbb; position: relative;
        }
        .acc-ui::before {
            content: ''; position: absolute; top: 2px; left: 2px;
            width: 3px; height: 3px; background: #ff4444; border-radius: 50%;
            box-shadow: 5px 0 0 #4488ff, 2px 5px 0 #44cc44, 6px 4px 0 #ffaa00;
        }
        /* Security: shield */
        .acc-sec {
            top: 26px; left: -6px; width: 12px; height: 14px;
            background: var(--color-security);
            clip-path: polygon(50% 0%, 100% 15%, 100% 55%, 50% 100%, 0% 55%, 0% 15%);
            border: 1px solid #d4830f;
        }
        .acc-sec::after {
            content: ''; position: absolute; top: 4px; left: 50%; transform: translateX(-50%);
            width: 4px; height: 5px; background: #fff; opacity: 0.5;
            clip-path: polygon(50% 0%, 100% 30%, 80% 100%, 20% 100%, 0% 30%);
        }
        /* Architect: ruler/blueprint roll */
        .acc-arch {
            top: 26px; right: -6px; width: 6px; height: 18px;
            background: linear-gradient(90deg, #4488cc, #5599dd);
            border: 1px solid #3377bb; border-radius: 2px;
        }
        /* Sr Dev: -> terminal prompt */
        /* (uses rubber duck on desk instead) */
        /* Jr Dev: notebook */
        .acc-jr {
            top: 28px; right: -4px; width: 10px; height: 12px;
            background: #ff9; border: 1px solid #cc9;
            border-left: 3px solid #c44;
        }

        /* ---- Agent Label (always visible) ---- */
        .agent-label {
            position: absolute; bottom: -18px; left: 50%; transform: translateX(-50%);
            font-size: 8px; letter-spacing: 1px; white-space: nowrap;
            padding: 2px 6px; text-transform: uppercase;
            border: 2px solid; transition: opacity 0.15s;
        }

        /* (status indicators replaced by nameplate dot) */

        /* ---- Speech Bubble ---- */
        .bubble {
            position: absolute; z-index: 42; padding: 4px 8px;
            background: rgba(20,14,30,0.95); border: 2px solid;
            font-size: 7px; max-width: 140px; word-wrap: break-word;
            font-family: 'Courier New', monospace;
            animation: bubIn .2s step-end, bubOut .5s ease-out 3.5s forwards;
            pointer-events: none;
        }
        .bubble::after {
            content: ''; position: absolute; bottom: -6px; left: 14px;
            border-left: 5px solid transparent; border-right: 5px solid transparent;
            border-top: 6px solid; border-top-color: inherit;
        }
        @keyframes bubIn { 0%{transform:scale(0)}100%{transform:scale(1)} }
        @keyframes bubOut { to{opacity:0} }

        /* ---- Paper Airplanes (arcing) ---- */
        .airplane {
            position: absolute; z-index: 45; pointer-events: none;
            width: 0; height: 0;
            border-left: 14px solid;
            border-top: 6px solid transparent;
            border-bottom: 6px solid transparent;
            filter: drop-shadow(0 2px 3px rgba(0,0,0,.4));
            animation: flyPlane var(--dur,1.4s) ease-in-out forwards;
        }
        @keyframes flyPlane {
            0%{left:var(--sx);top:var(--sy);opacity:1;transform:rotate(var(--sa,0deg))}
            40%{top:var(--arc-y);opacity:1;transform:rotate(var(--ma,-10deg))}
            100%{left:var(--ex);top:var(--ey);opacity:0;transform:rotate(var(--ea,0deg))}
        }

        /* Landing sparkle */
        .land-sparkle {
            position: absolute; z-index: 46; pointer-events: none;
            width: 16px; height: 16px;
            animation: sparkleOut 0.6s ease-out forwards;
        }
        .land-sparkle::before, .land-sparkle::after {
            content: '+'; position: absolute; font-size: 12px; font-weight: bold;
        }
        .land-sparkle::before { top: 0; left: 4px; }
        .land-sparkle::after { top: 4px; left: 0; font-size: 8px; }
        @keyframes sparkleOut {
            0%{transform:scale(0);opacity:1}
            50%{transform:scale(1.5);opacity:1}
            100%{transform:scale(0.5);opacity:0}
        }

        /* Desk lamp glow */
        .desk-lamp-glow {
            position: absolute; z-index: 3; pointer-events: none;
            width: 60px; height: 30px; border-radius: 50%;
            background: radial-gradient(ellipse, rgba(255,240,200,0.12) 0%, transparent 70%);
        }
        .desk-lamp-glow.on {
            background: radial-gradient(ellipse, rgba(51,255,51,0.1) 0%, transparent 70%);
        }

        /* ---- Post-it Notes ---- */
        .postit {
            position: absolute; width: 28px; height: 28px; z-index: 15;
            font-size: 5px; padding: 3px; line-height: 1.3; overflow: hidden;
            font-family: 'Courier New', monospace;
            box-shadow: 2px 2px 4px rgba(0,0,0,.3);
            transform: rotate(var(--rot,-2deg));
            cursor: pointer; transition: transform .2s, z-index 0s;
        }
        .postit:hover { transform: scale(2.5) rotate(0deg); z-index: 60; }

        /* ---- Atmosphere ---- */
        .atmo {
            position: absolute; top: 0; left: 0; right: 0; bottom: 0;
            pointer-events: none; z-index: 50;
            transition: background 1s, opacity 1s;
            opacity: 0;
        }

        /* Pizza boxes (appear during dev phase) */
        .pizza {
            position: absolute; z-index: 6; display: none;
        }
        .pizza-box {
            width: 24px; height: 20px; background: #c8a060;
            border: 2px solid #a08040; position: relative;
        }
        .pizza-box::before {
            content: ''; position: absolute; top: 2px; left: 50%; transform: translateX(-50%);
            width: 12px; height: 12px; background: #e8c880;
            border-radius: 50%; border: 1px solid #c8a060;
        }
        .pizza-box::after {
            content: ''; position: absolute; top: 5px; left: 50%; transform: translateX(-50%);
            width: 8px; height: 8px; background: #cc4444;
            border-radius: 50%;
        }
        .phase-dev .pizza { display: block; }

        /* ======================== SIDE PANEL ======================== */
        .side-panel {
            width: 0; background: var(--bg-panel); border-left: 4px solid var(--border-color);
            overflow: hidden; transition: width .3s; display: flex; flex-direction: column;
            image-rendering: auto;
        }
        .side-panel.open { width: 300px; }
        .sp-head {
            display: flex; align-items: center; justify-content: space-between;
            padding: 10px 12px; border-bottom: 3px solid var(--border-color);
        }
        .sp-info { display: flex; align-items: center; gap: 8px; }
        .sp-icon {
            width: 28px; height: 28px; display: flex; align-items: center;
            justify-content: center; font-size: 12px; border: 3px solid;
        }
        .sp-name { font-size: 8px; letter-spacing: 1px; }
        .sp-role { font-size: 6px; color: var(--text-secondary); margin-top: 2px; }
        .sp-close {
            background: none; border: 2px solid var(--border-color); color: var(--text-secondary);
            width: 24px; height: 24px; cursor: pointer; font-family: inherit; font-size: 10px;
            display: flex; align-items: center; justify-content: center;
        }
        .sp-close:hover { border-color: #E74C3C; color: #E74C3C; }
        .sp-body { flex: 1; overflow-y: auto; padding: 10px 12px; }
        .sp-section { margin-bottom: 14px; }
        .sp-title {
            font-size: 6px; letter-spacing: 2px; color: var(--text-dim);
            margin-bottom: 6px; padding-bottom: 4px; border-bottom: 1px solid var(--border-color);
            text-transform: uppercase;
        }
        .sp-task {
            padding: 6px 8px; background: var(--bg-dark); border: 2px solid var(--border-color);
            margin-bottom: 4px;
        }
        .sp-task-title { font-size: 7px; margin-bottom: 4px; line-height: 1.4; }
        .sp-bar { height: 6px; background: var(--bg-dark); border: 2px solid var(--border-color); overflow: hidden; }
        .sp-fill {
            height: 100%; transition: width .5s; position: relative;
        }
        .sp-pct { font-size: 6px; color: var(--text-secondary); margin-top: 2px; text-align: right; }
        .sp-msg {
            padding: 4px 6px; border-left: 3px solid var(--border-color);
            margin-bottom: 4px; font-family: 'Courier New', monospace;
        }
        .sp-msg-from { font-size: 6px; letter-spacing: 1px; margin-bottom: 1px; }
        .sp-msg-text { font-size: 7px; color: var(--text-secondary); max-height: 36px; overflow: hidden; line-height: 1.4; }
        .sp-msg-time { font-size: 5px; color: var(--text-dim); margin-top: 1px; }
        .sp-status-row {
            display: flex; align-items: center; gap: 8px; font-size: 8px;
        }
        .sp-status-dot { width: 8px; height: 8px; }

        /* ======================== BOTTOM BAR ======================== */
        .bottom-bar {
            background: var(--bg-dark); border-top: 4px solid var(--border-color);
            display: flex; flex-direction: column; max-height: 180px; flex-shrink: 0;
        }
        .input-row {
            display: flex; padding: 6px 12px; gap: 8px;
            border-bottom: 2px solid var(--border-color);
        }
        .task-inp {
            flex: 1; background: var(--bg-panel); border: 2px solid var(--border-color);
            color: var(--text-primary); font-family: inherit; font-size: 8px;
            padding: 6px 10px; outline: none;
        }
        .task-inp:focus { border-color: var(--accent); }
        .task-inp::placeholder { color: var(--text-dim); }
        .task-btn {
            padding: 6px 16px; background: var(--accent); border: 3px solid #3aaa9a;
            color: var(--bg-dark); font-family: inherit; font-size: 8px;
            cursor: pointer; letter-spacing: 1px;
        }
        .task-btn:hover { background: #5de0d6; }
        .task-btn:disabled { background: #444; border-color: #333; cursor: not-allowed; }
        .msg-log {
            flex: 1; overflow-y: auto; padding: 4px 12px; max-height: 100px;
        }
        .log-row {
            padding: 2px 0; display: flex; gap: 8px; font-family: 'Courier New', monospace;
            border-bottom: 1px solid rgba(58,40,80,0.4); font-size: 8px;
        }
        .log-time { color: var(--text-dim); font-size: 7px; min-width: 55px; }
        .log-who { min-width: 50px; font-size: 7px; }
        .log-text { color: var(--text-secondary); flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

        /* ======================== SCROLLBAR ======================== */
        ::-webkit-scrollbar{width:6px}
        ::-webkit-scrollbar-track{background:var(--bg-dark)}
        ::-webkit-scrollbar-thumb{background:var(--border-color)}

        /* ======================== RESPONSIVE ======================== */
        @media(max-width:1100px){
            .office-scene{transform:translate(-50%,-50%) scale(.75)}
            .side-panel.open{width:260px}
        }
        /* Wall Clock */
        .wall-clock {
            position: absolute; top: 10px; right: 130px;
            width: 48px; height: 48px; background: #443060;
            border: 4px solid #665090; z-index: 5;
        }
        .clock-inner {
            width: 100%; height: 100%; position: relative;
            display: flex; align-items: center; justify-content: center;
            background: #1a1030;
        }
        .clock-h, .clock-m {
            position: absolute; bottom: 50%; left: 50%;
            transform-origin: bottom center; background: var(--accent);
        }
        .clock-h { width: 3px; height: 12px; }
        .clock-m { width: 2px; height: 16px; }
        .clock-dot {
            width: 4px; height: 4px; background: var(--accent); position: absolute;
        }

        /* Dust Particles */
        .dust {
            position: absolute; width: 2px; height: 2px; background: rgba(200,180,160,0.15);
            border-radius: 50%; pointer-events: none; z-index: 3;
            animation: dustFloat var(--dust-dur, 12s) linear infinite;
        }
        @keyframes dustFloat {
            0% { transform: translate(0, 0); opacity: 0; }
            10% { opacity: 0.3; }
            90% { opacity: 0.2; }
            100% { transform: translate(var(--dust-dx, 40px), var(--dust-dy, -60px)); opacity: 0; }
        }

        /* Phase Banner */
        .phase-banner {
            position: absolute; top: 0; left: 0; right: 0; z-index: 55;
            text-align: center; padding: 8px; font-size: 10px;
            letter-spacing: 3px; transform: translateY(-100%);
            transition: transform 0.4s; pointer-events: none;
        }
        .phase-banner.show {
            transform: translateY(0);
        }

        /* Desk wrapper for click targets */
        .desk-unit {
            position: absolute; z-index: 5; cursor: pointer;
        }
        .desk-unit:hover .desk-glow {
            opacity: 1;
        }
        .desk-glow {
            position: absolute; top: -30px; left: -8px; right: -8px; bottom: -80px;
            border: 2px solid; border-radius: 4px; opacity: 0;
            transition: opacity 0.2s; pointer-events: none;
        }

        /* Idle behaviors */
        .idle-look-left .head { transform: translateX(-2px); }
        .idle-look-right .head { transform: translateX(2px); }
        .idle-lean .torso { transform: rotate(-3deg); }
        .idle-stretch .arm.la { transform: translateY(-6px) rotate(-15deg); }
        .idle-stretch .arm.ra { transform: translateY(-6px) rotate(15deg); }
        .idle-scan .eye.l { animation: scanEye 1s step-end; }
        .idle-scan .eye.r { animation: scanEye 1s step-end; }
        @keyframes scanEye { 25%{transform:translateX(-1px)} 75%{transform:translateX(1px)} }
        .idle-peek { transform: translateX(-12px); transition: transform 0.4s; }
        .idle-tap .acc-pm { animation: tapClip 0.5s step-end 2; }
        @keyframes tapClip { 50%{transform:rotate(5deg)} }
        .idle-nod .head { animation: nodHead 0.8s ease-in-out; }
        @keyframes nodHead { 30%{transform:translateY(2px)} 60%{transform:translateY(-1px)} }

        /* Cat roaming */
        .pixel-cat.roaming {
            transition: left 3s ease-in-out, top 2s ease-in-out;
        }
        .pixel-cat.awake .cat-zzz { display: none; }
        .pixel-cat.awake .cat-body {
            animation: catWake 0.5s ease-out;
        }
        @keyframes catWake {
            0%{transform:scaleY(0.6) scaleX(1.3)}
            50%{transform:scaleY(1.2) scaleX(0.9)}
            100%{transform:scaleY(1) scaleX(1)}
        }

        /* Coffee machine brewing */
        .coffee-machine.brewing .cm-screen {
            background: #0a5a0a;
            animation: brewGlow 0.5s step-end 4;
        }
        @keyframes brewGlow { 50%{background:#0a8a0a} }
        .coffee-machine .brew-steam {
            position: absolute; top: -8px; left: 50%; transform: translateX(-50%);
            font-size: 10px; color: rgba(255,255,255,0.3); display: none;
        }
        .coffee-machine.brewing .brew-steam { display: block; animation: steam 1s ease-in-out infinite; }
        .coffee-machine .brew-count {
            font-size: 6px; text-align: center; color: var(--accent);
            margin-top: 2px; letter-spacing: 1px;
        }

        /* Whiteboard modal */
        .wb-modal {
            position: fixed; top: 0; left: 0; right: 0; bottom: 0;
            background: rgba(10,6,14,0.85); z-index: 200;
            display: flex; align-items: center; justify-content: center;
            animation: fadeIn 0.2s;
        }
        @keyframes fadeIn { 0%{opacity:0}100%{opacity:1} }
        .wb-modal-content {
            width: 400px; max-height: 500px; background: #eee8dd;
            border: 6px solid #888; padding: 0; position: relative;
            image-rendering: auto;
        }
        .wb-modal-header {
            background: #cc4444; color: #fff; font-size: 10px;
            text-align: center; padding: 8px; letter-spacing: 2px;
        }
        .wb-modal-body {
            padding: 16px; font-size: 9px; color: #333;
            font-family: 'Courier New', monospace; line-height: 1.8;
            max-height: 400px; overflow-y: auto;
        }
        .wb-modal-body .wb-entry {
            padding: 4px 0; border-bottom: 1px dashed #ccc;
        }
        .wb-modal-body .wb-entry::before { content: '> '; color: #999; }
        .wb-modal-close {
            position: absolute; top: 6px; right: 8px; background: none;
            border: none; color: #fff; font-size: 14px; cursor: pointer;
            font-family: inherit;
        }

        /* Day/night tint on windows */
        .window.night { background: #060e1a; }
        .window.dawn { background: linear-gradient(180deg, #1a2040 0%, #804020 100%); }
        .window.day { background: linear-gradient(180deg, #4080c0 0%, #80b0e0 100%); }
        .window.dusk { background: linear-gradient(180deg, #2a1840 0%, #c05030 100%); }
        .window.day .star { display: none; }
        .window.dawn .star { opacity: 0.3; }
        .window.dusk .star { opacity: 0.5; }

        /* ======================== SPRINT 3: DELIGHT ======================== */

        /* XP Bar in panel */
        .xp-section { margin-bottom: 14px; }
        .xp-row {
            display: flex; align-items: center; gap: 8px; margin-bottom: 6px;
        }
        .xp-level {
            font-size: 8px; letter-spacing: 1px; min-width: 48px;
        }
        .xp-bar-wrap {
            flex: 1; height: 8px; background: var(--bg-dark);
            border: 2px solid var(--border-color); overflow: hidden;
        }
        .xp-bar-fill {
            height: 100%; transition: width 0.6s;
        }
        .xp-text { font-size: 6px; color: var(--text-dim); text-align: right; }

        /* Level up animation */
        .level-up-pop {
            position: absolute; z-index: 60; pointer-events: none;
            font-size: 10px; letter-spacing: 2px; color: #FFD700;
            text-shadow: 0 0 8px rgba(255,215,0,0.6);
            animation: lvlUp 2s ease-out forwards;
        }
        @keyframes lvlUp {
            0%{transform:translateX(-50%) translateY(0) scale(0.5);opacity:0}
            20%{transform:translateX(-50%) translateY(-10px) scale(1.2);opacity:1}
            60%{transform:translateX(-50%) translateY(-25px) scale(1);opacity:1}
            100%{transform:translateX(-50%) translateY(-40px) scale(0.8);opacity:0}
        }

        /* Achievement toast */
        .achievement-toast {
            position: fixed; top: 60px; right: 20px; z-index: 300;
            background: var(--bg-panel); border: 3px solid #FFD700;
            padding: 10px 16px; display: flex; align-items: center; gap: 10px;
            animation: toastIn 0.3s ease-out, toastOut 0.4s ease-in 3.6s forwards;
            image-rendering: auto; max-width: 320px;
            box-shadow: 0 4px 20px rgba(255,215,0,0.2);
        }
        @keyframes toastIn {
            0%{transform:translateX(120%);opacity:0}
            100%{transform:translateX(0);opacity:1}
        }
        @keyframes toastOut {
            to{transform:translateX(120%);opacity:0}
        }
        .achievement-icon {
            font-size: 20px; min-width: 28px; text-align: center;
        }
        .achievement-info { flex: 1; }
        .achievement-label {
            font-size: 6px; color: #FFD700; letter-spacing: 2px;
            text-transform: uppercase; margin-bottom: 2px;
        }
        .achievement-name { font-size: 8px; color: var(--text-primary); }
        .achievement-desc { font-size: 6px; color: var(--text-secondary); margin-top: 2px; }

        /* Stats bar in header */
        .stats-bar {
            display: flex; gap: 12px; align-items: center;
            font-size: 7px; color: var(--text-dim);
        }
        .stat-item { display: flex; align-items: center; gap: 4px; }
        .stat-val { color: var(--accent); }

        /* Trophy shelf in panel */
        .trophy-shelf {
            display: flex; flex-wrap: wrap; gap: 4px; padding: 4px 0;
        }
        .trophy {
            width: 24px; height: 24px; display: flex; align-items: center;
            justify-content: center; font-size: 14px; background: var(--bg-dark);
            border: 2px solid var(--border-color); cursor: pointer;
            position: relative;
        }
        .trophy.locked { opacity: 0.3; filter: grayscale(1); }
        .trophy:hover::after {
            content: attr(data-name); position: absolute; bottom: -18px;
            left: 50%; transform: translateX(-50%); white-space: nowrap;
            font-size: 5px; background: var(--bg-dark); border: 1px solid var(--border-color);
            padding: 2px 4px; z-index: 10; color: var(--text-secondary);
        }

        /* Crunch mode */
        .crunch-mode .office-scene { animation: crunchShake 0.1s step-end infinite; }
        @keyframes crunchShake { 50%{transform:translate(-50%,-50%) translate(1px,-1px)} }
        .crunch-mode .monitor-screen { background: #2a0a0a !important; }
        .crunch-mode .monitor-screen.on {
            background: #3a0a0a !important;
            animation: crunchGlow 0.5s step-end infinite !important;
        }
        @keyframes crunchGlow { 50%{box-shadow:inset 0 0 15px rgba(255,50,50,0.5)} }
        .crunch-mode .pizza { display: block !important; }
        .crunch-banner {
            position: fixed; top: 50%; left: 50%; transform: translate(-50%,-50%);
            z-index: 250; font-size: 16px; color: #ff4444; letter-spacing: 4px;
            text-shadow: 0 0 20px rgba(255,0,0,0.5);
            animation: crunchText 0.5s step-end 3;
            pointer-events: none;
        }
        @keyframes crunchText {
            0%{transform:translate(-50%,-50%) scale(3);opacity:0}
            50%{transform:translate(-50%,-50%) scale(1);opacity:1}
            100%{opacity:0}
        }

        @media(max-width:800px){
            .office-scene{transform:translate(-50%,-50%) scale(.55)}
            .side-panel{position:absolute;right:0;top:44px;bottom:0;z-index:100}
        }
    </style>
</head>
<body>
<div class="app">
    <!-- Header -->
    <div class="header">
        <div class="header-left">
            <div class="logo">AGENT HOUSE</div>
            <div class="phase-pill"><div class="phase-dot"></div><span id="phaseText">IDLE</span></div>
            <div class="stats-bar" id="statsBar">
                <div class="stat-item">LVL <span class="stat-val" id="statLevel">1</span></div>
                <div class="stat-item">BUILDS <span class="stat-val" id="statBuilds">0</span></div>
                <div class="stat-item">FILES <span class="stat-val" id="statFiles">0</span></div>
            </div>
        </div>
        <div class="header-right">
            <select class="project-sel" id="projectSelect"><option value="default">default</option></select>
            <button class="new-proj-btn" id="newProjBtn" title="Create new project">+</button>
            <a class="view-toggle" href="/">DASHBOARD</a>
            <div class="ws-dot" id="wsDot" title="Disconnected"></div>
        </div>
    </div>

    <div class="main-area">
        <!-- Office -->
        <div class="office-wrap">
            <div class="office-scene" id="scene">
                <!-- Wall -->
                <div class="wall">
                    <div class="window w1">
                        <div class="window-pane"></div><div class="window-pane h"></div>
                        <div class="star" style="top:12px;left:10px"></div>
                        <div class="star" style="top:28px;left:45px;animation-delay:.8s"></div>
                        <div class="star" style="top:50px;left:20px;animation-delay:1.5s"></div>
                    </div>
                    <div class="window w2">
                        <div class="window-pane"></div><div class="window-pane h"></div>
                        <div class="star" style="top:15px;left:50px;animation-delay:.3s"></div>
                        <div class="star" style="top:40px;left:15px;animation-delay:1.1s"></div>
                    </div>
                    <div class="window w3">
                        <div class="window-pane"></div><div class="window-pane h"></div>
                        <div class="star" style="top:10px;left:30px;animation-delay:.6s"></div>
                        <div class="star" style="top:45px;left:55px;animation-delay:1.8s"></div>
                    </div>
                    <div class="window w4">
                        <div class="window-pane"></div><div class="window-pane h"></div>
                        <div class="star" style="top:20px;left:15px;animation-delay:.4s"></div>
                        <div class="star" style="top:50px;left:48px;animation-delay:1.3s"></div>
                    </div>

                    <!-- Posters -->
                    <div class="poster p1">
                        <span class="poster-art">&#9734;</span>
                        <span class="poster-text">Ship It!</span>
                    </div>
                    <div class="poster p2">
                        <span class="poster-art">&#9888;</span>
                        <span class="poster-text">OWASP</span>
                    </div>
                    <div class="poster p3">
                        <span class="poster-art">&#9829;</span>
                        <span class="poster-text">UX First</span>
                    </div>

                    <!-- Clock -->
                    <div class="wall-clock">
                        <div class="clock-inner">
                            <div class="clock-dot"></div>
                            <div class="clock-h" id="clockH"></div>
                            <div class="clock-m" id="clockM"></div>
                        </div>
                    </div>

                    <!-- Whiteboard -->
                    <div class="whiteboard" id="whiteboard">
                        <div class="wb-header">SPRINT BOARD</div>
                        <div class="wb-body" id="wbBody">
                            <span class="wb-line">Waiting for task...</span>
                        </div>
                    </div>
                </div>

                <!-- Floor -->
                <div class="floor" id="floorArea">
                    <!-- Zone carpets -->
                    <div class="zone exec"><div class="zone-label">Executive</div></div>
                    <div class="zone design"><div class="zone-label">Design Pod</div></div>
                    <div class="zone tech"><div class="zone-label">Tech Corner</div></div>
                    <div class="zone dev"><div class="zone-label">Dev Pit</div></div>

                    <!-- Bookshelves on wall -->
                    <div class="bookshelf" style="left:12px;top:8px">
                        <div class="shelf-frame">
                            <div class="shelf-row">
                                <div class="book" style="height:12px;background:#c44"></div>
                                <div class="book" style="height:10px;background:#48c"></div>
                                <div class="book" style="height:14px;background:#4a4"></div>
                                <div class="book" style="height:11px;background:#c84"></div>
                            </div>
                            <div class="shelf-row">
                                <div class="book" style="height:13px;background:#84c"></div>
                                <div class="book" style="height:9px;background:#cc4"></div>
                                <div class="book" style="height:12px;background:#4cc"></div>
                            </div>
                        </div>
                    </div>

                    <!-- Bookshelf 2 (tech corner) -->
                    <div class="bookshelf" style="right:12px;top:8px">
                        <div class="shelf-frame">
                            <div class="shelf-row">
                                <div class="book" style="height:11px;background:#94c"></div>
                                <div class="book" style="height:13px;background:#4a8"></div>
                                <div class="book" style="height:10px;background:#c64"></div>
                                <div class="book" style="height:14px;background:#48c"></div>
                            </div>
                            <div class="shelf-row">
                                <div class="book" style="height:12px;background:#ca4"></div>
                                <div class="book" style="height:10px;background:#4cc"></div>
                                <div class="book" style="height:11px;background:#c44"></div>
                            </div>
                        </div>
                    </div>

                    <!-- Plants -->
                    <div class="plant" style="right:60px;top:8px"><div class="plant-top"></div><div class="plant-pot"></div></div>
                    <div class="plant" style="left:10px;top:220px"><div class="plant-top"></div><div class="plant-pot"></div></div>
                    <div class="plant" style="right:10px;top:240px"><div class="plant-top"></div><div class="plant-pot"></div></div>

                    <!-- Coffee Machine -->
                    <div class="coffee-machine" style="left:580px;top:130px" id="coffeeMachine">
                        <div class="cm-top"></div>
                        <div class="cm-body">
                            <div class="cm-screen"></div>
                            <div class="cm-spout"></div>
                        </div>
                        <div class="cm-cup"></div>
                        <div class="brew-steam">~</div>
                        <div class="brew-count" id="brewCount"></div>
                    </div>

                    <!-- Cat on UX desk -->
                    <div class="pixel-cat" id="officeCat" style="left:560px;top:50px" title="Office cat - Mr. Whiskers">
                        <div class="cat-zzz">z</div>
                        <div class="cat-body"></div>
                        <div class="cat-tail"></div>
                    </div>

                    <!-- Pizza boxes (shown during dev phase) -->
                    <div class="pizza" id="pizza1" style="left:460px;top:280px"><div class="pizza-box"></div></div>
                    <div class="pizza" id="pizza2" style="left:580px;top:290px"><div class="pizza-box"></div></div>
                </div>

                <!-- Agent layer (rendered by JS) -->
                <div id="agentLayer"></div>

                <!-- Phase banner -->
                <div class="phase-banner" id="phaseBanner"></div>

                <!-- Dust particles -->
                <div class="dust" style="left:100px;top:200px;--dust-dur:15s;--dust-dx:30px;--dust-dy:-50px"></div>
                <div class="dust" style="left:300px;top:280px;--dust-dur:18s;--dust-dx:-20px;--dust-dy:-40px;animation-delay:3s"></div>
                <div class="dust" style="left:600px;top:220px;--dust-dur:14s;--dust-dx:25px;--dust-dy:-55px;animation-delay:6s"></div>
                <div class="dust" style="left:800px;top:300px;--dust-dur:16s;--dust-dx:-30px;--dust-dy:-45px;animation-delay:9s"></div>
                <div class="dust" style="left:450px;top:350px;--dust-dur:20s;--dust-dx:15px;--dust-dy:-60px;animation-delay:2s"></div>
                <div class="dust" style="left:200px;top:320px;--dust-dur:17s;--dust-dx:-25px;--dust-dy:-50px;animation-delay:7s"></div>

                <!-- Atmosphere overlay -->
                <div class="atmo" id="atmo"></div>
            </div>
        </div>

        <!-- Side Panel -->
        <div class="side-panel" id="sidePanel">
            <div class="sp-head">
                <div class="sp-info">
                    <div class="sp-icon" id="spIcon"></div>
                    <div><div class="sp-name" id="spName"></div><div class="sp-role" id="spRole"></div></div>
                </div>
                <button class="sp-close" id="spClose">x</button>
            </div>
            <div class="sp-body" id="spBody"></div>
        </div>
    </div>

    <!-- Bottom Bar -->
    <div class="bottom-bar">
        <div class="input-row">
            <input type="text" class="task-inp" id="taskInp" placeholder="Describe your project...">
            <button class="task-btn" id="taskBtn">BUILD</button>
        </div>
        <div class="msg-log" id="msgLog"></div>
    </div>
</div>

<script>
// ============================
// LOGGER
// ============================
const LOG_PREFIX = '[OFFICE]';
const L = {
    info:  (...a) => console.log(LOG_PREFIX, ...a),
    warn:  (...a) => console.warn(LOG_PREFIX, ...a),
    error: (...a) => console.error(LOG_PREFIX, ...a),
    debug: (...a) => console.debug(LOG_PREFIX, ...a),
    ws:    (...a) => console.log('[WS]', ...a),
    api:   (...a) => console.log('[API]', ...a),
    agent: (...a) => console.log('[AGENT]', ...a),
    phase: (...a) => console.log('[PHASE]', ...a),
    task:  (...a) => console.log('[TASK]', ...a),
    xp:    (...a) => console.log('[XP]', ...a),
};

// ============================
// STATE
// ============================
const S = {
    ws: null,
    agents: {},
    messages: [],
    projectId: 'default',
    phase: '',
    running: false,
    selected: null,
    wbNotes: []
};

// Agent defs with positions (relative to floor, which starts at y=160 in scene)
// So deskY here is relative to the scene top
const AGENTS = {
    ceo: {
        role:'ceo', name:'CEO', full:'Chief Executive Officer',
        color:'#4A90A4', icon:'C', skin:'#FDDCB5',
        hair:'ceo', acc:'ceo',
        deskX:60, deskY:180,
        screenIdle:'VISION\\nSTRATEGY\\nGROWTH',
        deskItems: []
    },
    pm: {
        role:'pm', name:'PM', full:'Product Manager',
        color:'#7B68EE', icon:'P', skin:'#D4A07A',
        hair:'pm', acc:'pm',
        deskX:170, deskY:180,
        screenIdle:'FEATURES\\nSCOPE\\nUSERS',
        deskItems: ['coffee']
    },
    ux: {
        role:'ux', name:'UX', full:'UX Designer',
        color:'#FF6B6B', icon:'X', skin:'#FDE0C8',
        hair:'ux', acc:'ux',
        deskX:330, deskY:180,
        screenIdle:'FLOWS\\nWIRES\\nTESTS',
        deskItems: []
    },
    ui: {
        role:'ui', name:'UI', full:'UI Designer',
        color:'#4ECDC4', icon:'U', skin:'#C68E5B',
        hair:'ui', acc:'ui',
        deskX:460, deskY:180,
        screenIdle:'COLORS\\nFONTS\\nGRID',
        deskItems: []
    },
    security: {
        role:'security', name:'SEC', full:'Security Expert',
        color:'#F39C12', icon:'S', skin:'#8D5524',
        hair:'sec', acc:'sec',
        deskX:640, deskY:180,
        screenIdle:'THREATS\\nAUDIT\\nOWASP',
        deskItems: ['coffee']
    },
    architect: {
        role:'architect', name:'ARCH', full:'Software Architect',
        color:'#9B59B6', icon:'A', skin:'#E8C498',
        hair:'arch', acc:'arch', hasGlasses: true,
        deskX:800, deskY:180,
        screenIdle:'SCHEMA\\nSTACK\\nAPI',
        deskItems: []
    },
    'senior-dev': {
        role:'senior-dev', name:'SR DEV', full:'Senior Developer',
        color:'#2ECC71', icon:'>', skin:'#F0C8A0',
        hair:'sr', acc:null, hasBeard: true,
        deskX:360, deskY:320,
        screenIdle:'func()\\nimport\\ntest()',
        deskItems: ['duck']
    },
    'junior-dev': {
        role:'junior-dev', name:'JR DEV', full:'Junior Developer',
        color:'#3498DB', icon:'_', skin:'#FDDCB5',
        hair:'jr', acc:'jr', hasCap: true,
        deskX:540, deskY:320,
        screenIdle:'learn()\\ncode()\\npush()',
        deskItems: ['coffee']
    }
};

const PHASE_ATMO = {
    research:          {color:'rgba(84,160,255,0.15)',  label:'RESEARCHING'},
    planning:          {color:'rgba(95,39,205,0.15)',   label:'PLANNING'},
    discussion:        {color:'rgba(0,210,211,0.15)',   label:'DISCUSSING'},
    development:       {color:'rgba(16,172,132,0.15)',  label:'BUILDING'},
    template_selection:{color:'rgba(255,159,67,0.15)',  label:'SELECTING TEMPLATE'},
    qa:                {color:'rgba(231,76,60,0.15)',   label:'QA REVIEW'}
};

// ============================
// RENDER
// ============================
function render() {
    const layer = document.getElementById('agentLayer');
    layer.innerHTML = '';

    Object.entries(AGENTS).forEach(([key, ag]) => {
        const as = S.agents[key] || {state:'idle', progress:0, taskTitle:''};
        const working = as.state === 'working';
        const error = as.state === 'error';
        const done = as.state === 'completed';
        const isLarge = key === 'ceo';

        // ---- Desk Unit (below character = closer to user) ----
        const unit = document.createElement('div');
        unit.className = 'desk-unit';
        unit.style.cssText = 'left:'+ag.deskX+'px;top:'+(ag.deskY+50)+'px';
        unit.addEventListener('click', () => openPanel(key));

        let deskItemsHTML = '';
        ag.deskItems.forEach((item, i) => {
            if (item === 'duck') {
                deskItemsHTML += '<div class="desk-item" style="position:absolute;top:4px;right:6px"><div class="rubber-duck"></div></div>';
            } else if (item === 'coffee') {
                deskItemsHTML += '<div class="desk-item" style="position:absolute;top:2px;right:6px"><div class="coffee-cup"></div></div>';
            }
        });

        // Progress bar on desk
        let progHTML = '';
        if (working && as.progress > 0) {
            progHTML = '<div style="position:absolute;bottom:4px;left:8px;right:8px;height:4px;background:rgba(0,0,0,.3);overflow:hidden"><div style="height:100%;width:'+as.progress+'%;background:'+ag.color+';transition:width .5s"></div></div>';
        }

        // Desk glow on hover + lamp glow (glow from monitor on desk surface)
        const glowHTML = '<div class="desk-glow" style="border-color:'+ag.color+'44"></div>';
        const lampHTML = '<div class="desk-lamp-glow'+(working?' on':'')+'" style="top:-4px;left:20px"></div>';

        // Status dot class
        const dotClass = working ? 'working' : error ? 'error' : done ? 'completed' : 'idle';

        // Desk with monitor-back on top edge, nameplate on bottom
        unit.innerHTML =
            glowHTML + lampHTML +
            '<div class="desk-top'+(isLarge?' large':'')+'">' +
                '<div class="monitor-back'+(working?' on':'')+'">' +
                    '<div class="monitor-back-panel"></div>' +
                '</div>' +
                '<div class="monitor-stand"></div>' +
                deskItemsHTML +
                progHTML +
            '</div>' +
            '<div class="nameplate-front" style="color:'+ag.color+';border-color:'+ag.color+';background:rgba(20,14,30,.95)">' +
                '<div class="status-dot '+dotClass+'" style="color:'+ag.color+'"></div>' +
                ag.name +
            '</div>';

        layer.appendChild(unit);

        // ---- Character (above desk = farther from user, facing forward) ----
        const el = document.createElement('div');
        el.className = 'agent';
        el.style.cssText = 'left:'+(ag.deskX+34)+'px;top:'+(ag.deskY-4)+'px';
        el.dataset.agent = key;

        // Accessories
        let accHTML = '';
        if (ag.acc === 'ceo') accHTML = '<div class="acc acc-ceo"></div>';
        else if (ag.acc === 'pm') accHTML = '<div class="acc acc-pm"></div>';
        else if (ag.acc === 'ux') accHTML = '<div class="acc acc-ux"></div>';
        else if (ag.acc === 'ui') accHTML = '<div class="acc acc-ui"></div>';
        else if (ag.acc === 'sec') accHTML = '<div class="acc acc-sec"></div>';
        else if (ag.acc === 'arch') accHTML = '<div class="acc acc-arch"></div>';
        else if (ag.acc === 'jr') accHTML = '<div class="acc acc-jr"></div>';

        let extraHTML = '';
        if (ag.hasGlasses) extraHTML += '<div class="glasses"></div>';
        if (ag.hasBeard) extraHTML += '<div class="beard"></div>';
        if (ag.hasCap) extraHTML += '<div class="cap-brim"></div>';

        el.innerHTML =
            '<div class="char '+(working?'typing':'')+'">' +
                '<div class="hair hair-'+ag.hair+'"></div>' +
                '<div class="head" style="background:'+ag.skin+';border-color:'+darken(ag.skin)+';">' +
                    '<div class="face'+(Math.random()>.75?' blink':'')+'">' +
                        '<div class="eye l"></div><div class="eye r"></div>' +
                        '<div class="mouth"></div>' +
                    '</div>' +
                '</div>' +
                extraHTML +
                accHTML +
                '<div class="torso" style="background:'+ag.color+';border-color:'+darken(ag.color)+';"></div>' +
                '<div class="arm la" style="background:'+ag.color+';border-color:'+darken(ag.color)+';"></div>' +
                '<div class="arm ra" style="background:'+ag.color+';border-color:'+darken(ag.color)+';"></div>' +
                '<div class="legs"><div class="leg"></div><div class="leg"></div></div>' +
            '</div>';

        el.addEventListener('click', () => openPanel(key));
        layer.appendChild(el);
    });
}

function darken(hex) {
    const r = Math.max(0, parseInt(hex.slice(1,3),16) - 30);
    const g = Math.max(0, parseInt(hex.slice(3,5),16) - 30);
    const b = Math.max(0, parseInt(hex.slice(5,7),16) - 30);
    return '#'+r.toString(16).padStart(2,'0')+g.toString(16).padStart(2,'0')+b.toString(16).padStart(2,'0');
}

// ============================
// PAPER AIRPLANE
// ============================
function sendAirplane(fromKey, toKey, color) {
    const f = AGENTS[fromKey], t = AGENTS[toKey];
    if (!f || !t) return;
    const scene = document.getElementById('scene');
    const el = document.createElement('div');
    el.className = 'airplane';
    const sx = f.deskX+50, sy = f.deskY+20, ex = t.deskX+50, ey = t.deskY+20;
    const arcY = Math.min(sy, ey) - 40 - Math.random()*30; // arc peak above both
    const angle = Math.atan2(ey-sy, ex-sx)*180/Math.PI;
    const dur = 1.0 + Math.random()*0.6;
    el.style.cssText =
        '--sx:'+sx+'px;--sy:'+sy+'px;--ex:'+ex+'px;--ey:'+ey+'px;'+
        '--arc-y:'+arcY+'px;'+
        '--sa:'+angle+'deg;--ma:'+(angle-20)+'deg;--ea:'+(angle+5)+'deg;'+
        '--dur:'+dur+'s;'+
        'border-left-color:'+(color||'#aaa')+';left:'+sx+'px;top:'+sy+'px;';
    scene.appendChild(el);
    // Landing sparkle
    setTimeout(() => {
        const sparkle = document.createElement('div');
        sparkle.className = 'land-sparkle';
        sparkle.style.cssText = 'left:'+ex+'px;top:'+ey+'px;color:'+(color||'#aaa')+';';
        scene.appendChild(sparkle);
        setTimeout(() => sparkle.remove(), 700);
    }, dur*900);
    setTimeout(() => el.remove(), (dur+0.5)*1000);
    playSound('whoosh');
}

// ============================
// SPEECH BUBBLE
// ============================
function showBubble(key, text) {
    const ag = AGENTS[key];
    if (!ag) return;
    const old = document.querySelector('.bubble[data-a="'+key+'"]');
    if (old) old.remove();
    const scene = document.getElementById('scene');
    const el = document.createElement('div');
    el.className = 'bubble';
    el.dataset.a = key;
    el.style.cssText = 'left:'+(ag.deskX+20)+'px;top:'+(ag.deskY+20)+'px;border-color:'+ag.color+';color:'+ag.color+';';

    // Typewriter effect
    const fullText = text.length > 45 ? text.substring(0,42)+'...' : text;
    el.textContent = '';
    scene.appendChild(el);
    let i = 0;
    const speed = key === 'ceo' ? 60 : key === 'junior-dev' ? 25 : 40; // CEO slow, Jr fast
    const typeInterval = setInterval(() => {
        if (i < fullText.length) {
            el.textContent += fullText[i];
            i++;
        } else {
            clearInterval(typeInterval);
        }
    }, speed);

    setTimeout(() => { clearInterval(typeInterval); el.remove(); }, 5000);
    playSound('blip');
}

// ============================
// POST-IT
// ============================
function addPostIt(key, text) {
    const ag = AGENTS[key];
    if (!ag) return;
    const scene = document.getElementById('scene');
    const old = scene.querySelectorAll('.postit[data-a="'+key+'"]');
    if (old.length > 2) old[0].remove();
    const el = document.createElement('div');
    el.className = 'postit';
    el.dataset.a = key;
    const ox = (Math.random()-.5)*20, oy = (Math.random()-.5)*10;
    el.style.cssText =
        'left:'+(ag.deskX+80+ox)+'px;top:'+(ag.deskY-8+oy)+'px;'+
        'background:'+ag.color+'25;border:1px solid '+ag.color+'55;color:'+ag.color+';'+
        '--rot:'+(Math.random()*10-5)+'deg;';
    el.textContent = text.length > 25 ? text.substring(0,22)+'...' : text;
    el.title = text;
    scene.appendChild(el);
}

// ============================
// WHITEBOARD
// ============================
function updateWB(note) {
    S.wbNotes.unshift(note);
    if (S.wbNotes.length > 5) S.wbNotes.pop();
    document.getElementById('wbBody').innerHTML = S.wbNotes
        .map(n => '<span class="wb-line">'+(n.length>22?n.substring(0,20)+'..':n)+'</span>').join('');
}

// ============================
// PHASE
// ============================
function setPhase(phase) {
    if (S.phase !== phase) L.phase('Phase change: '+(S.phase||'none')+' -> '+phase);
    S.phase = phase;
    const a = PHASE_ATMO[phase];
    const el = document.getElementById('atmo');
    const txt = document.getElementById('phaseText');
    const scene = document.getElementById('scene');

    if (a) {
        el.style.background = a.color;
        el.style.opacity = '1';
        txt.textContent = a.label;
        txt.style.color = 'var(--accent)';
        showPhaseBanner(a.label, a.color.replace('0.15','0.85'));
        playSound('phase');
    } else {
        el.style.opacity = '0';
        txt.textContent = S.running ? 'WORKING' : 'IDLE';
        txt.style.color = '';
    }

    // Pizza during dev
    if (phase === 'development') {
        scene.classList.add('phase-dev');
    } else {
        scene.classList.remove('phase-dev');
    }
}

// ============================
// PANEL
// ============================
function openPanel(key) {
    const ag = AGENTS[key];
    const as = S.agents[key] || {state:'idle',progress:0,tasks:[],messages:[]};
    S.selected = key;
    document.getElementById('sidePanel').classList.add('open');
    const icon = document.getElementById('spIcon');
    icon.textContent = ag.icon;
    icon.style.borderColor = ag.color;
    icon.style.color = ag.color;
    document.getElementById('spName').textContent = ag.name;
    document.getElementById('spName').style.color = ag.color;
    document.getElementById('spRole').textContent = ag.full;

    const stateLabel = as.state || 'idle';
    const stateColors = {idle:'#555',working:'#2ECC71',error:'#E74C3C',completed:'#4ECDC4',waiting:'#F39C12'};

    let tasksHTML = '';
    if (as.tasks && as.tasks.length > 0) {
        tasksHTML = as.tasks.map(t =>
            '<div class="sp-task"><div class="sp-task-title">'+esc(t.title||t.task_type||'Task')+'</div>'+
            '<div class="sp-bar"><div class="sp-fill" style="width:'+(t.progress||0)+'%;background:'+ag.color+'"></div></div>'+
            '<div class="sp-pct">'+(t.progress||0)+'% - '+(t.status||'pending')+'</div></div>'
        ).join('');
    } else if (as.taskTitle) {
        tasksHTML = '<div class="sp-task"><div class="sp-task-title">'+esc(as.taskTitle)+'</div>'+
            '<div class="sp-bar"><div class="sp-fill" style="width:'+(as.progress||0)+'%;background:'+ag.color+'"></div></div>'+
            '<div class="sp-pct">'+(as.progress||0)+'%</div></div>';
    } else {
        tasksHTML = '<div style="font-size:7px;color:var(--text-dim)">No active tasks</div>';
    }

    const agMsgs = S.messages.filter(m => m.from===key || m.to===key || m.from===ag.role).slice(-8);
    let msgsHTML = agMsgs.length ? agMsgs.map(m =>
        '<div class="sp-msg" style="border-left-color:'+(AGENTS[m.from]?.color||'#555')+'">'+
        '<div class="sp-msg-from" style="color:'+(AGENTS[m.from]?.color||'#888')+'">'+(m.from||'sys')+'</div>'+
        '<div class="sp-msg-text">'+esc((m.content||'').substring(0,100))+'</div>'+
        '<div class="sp-msg-time">'+fmtTime(m.timestamp)+'</div></div>'
    ).join('') : '<div style="font-size:7px;color:var(--text-dim)">No messages yet</div>';

    // XP section
    const xp = getXP(key);
    const xpInLevel = xp.xp % XP_PER_LEVEL;
    const xpPct = Math.min(100, (xpInLevel / XP_PER_LEVEL) * 100);
    const xpHTML =
        '<div class="xp-section">' +
        '<div class="xp-row">' +
        '<div class="xp-level" style="color:'+ag.color+'">LVL '+xp.level+'</div>' +
        '<div class="xp-bar-wrap"><div class="xp-bar-fill" style="width:'+xpPct+'%;background:'+ag.color+'"></div></div>' +
        '</div>' +
        '<div class="xp-text">'+xp.xp+' / '+(xp.level * XP_PER_LEVEL)+' XP</div>' +
        '</div>';

    // Trophy shelf
    let trophyHTML = '<div class="trophy-shelf">';
    Object.entries(ACHIEVEMENTS).forEach(([id, ach]) => {
        const unlocked = unlockedAchievements.includes(id);
        trophyHTML += '<div class="trophy'+(unlocked?'':' locked')+'" data-name="'+ach.name+'">'+ach.icon+'</div>';
    });
    trophyHTML += '</div>';

    document.getElementById('spBody').innerHTML =
        '<div class="sp-section"><div class="sp-title">Status</div>'+
        '<div class="sp-status-row"><div class="sp-status-dot" style="background:'+(stateColors[stateLabel]||'#555')+'"></div>'+
        '<span style="text-transform:uppercase;letter-spacing:1px">'+stateLabel+'</span></div></div>'+
        xpHTML +
        '<div class="sp-section"><div class="sp-title">Tasks</div>'+tasksHTML+'</div>'+
        '<div class="sp-section"><div class="sp-title">Activity</div>'+msgsHTML+'</div>'+
        '<div class="sp-section"><div class="sp-title">Achievements</div>'+trophyHTML+'</div>';
}

function closePanel() {
    document.getElementById('sidePanel').classList.remove('open');
    S.selected = null;
}

// ============================
// LOG
// ============================
function addLog(msg) {
    const log = document.getElementById('msgLog');
    const el = document.createElement('div');
    el.className = 'log-row';
    const c = AGENTS[msg.from]?.color || '#888';
    el.innerHTML = '<span class="log-time">'+fmtTime(msg.timestamp)+'</span>'+
        '<span class="log-who" style="color:'+c+'">'+(msg.from||'SYS')+'</span>'+
        '<span class="log-text">'+esc((msg.content||'').substring(0,200))+'</span>';
    log.appendChild(el);
    log.scrollTop = log.scrollHeight;
    while (log.children.length > 80) log.removeChild(log.firstChild);
}

// ============================
// WEBSOCKET
// ============================
function connectWS() {
    const proto = location.protocol==='https:'?'wss:':'ws:';
    const url = proto+'//'+location.host+'/ws';
    L.ws('Connecting to', url);
    S.ws = new WebSocket(url);
    S.ws.onopen = () => { L.ws('Connected'); document.getElementById('wsDot').classList.add('on'); };
    S.ws.onclose = (e) => { L.ws('Disconnected, code='+e.code+' reason='+e.reason+', reconnecting in 3s'); document.getElementById('wsDot').classList.remove('on'); setTimeout(connectWS,3000); };
    S.ws.onerror = (e) => { L.error('WebSocket error', e); S.ws.close(); };
    S.ws.onmessage = (e) => { try { handleWS(JSON.parse(e.data)); } catch(x){ L.error('WS parse error', x, 'raw:', e.data?.substring?.(0,200)); } };
}

function handleWS(data) {
    L.ws('Received:', data.type, data.message?.type||data.event?.event_type||'');
    if (data.type==='message' && data.message) {
        const m = data.message;
        L.debug('Message from='+m.from+' to='+(m.to||'-')+' type='+m.type+' content='+(m.content||'').substring(0,80));
        S.messages.push(m);
        if (S.messages.length>200) S.messages.shift();
        addLog(m);
        handleMsg(m);
    }
    if (data.type==='agent_task_event' && data.event) { L.task('Event:', data.event.event_type, 'agent='+data.event.agent_role, 'task='+data.event.task_id); handleTaskEvt(data.event); }
}

function handleMsg(msg) {
    const fk = normRole(msg.from), tk = normRole(msg.to);
    L.agent('handleMsg: from='+fk+' to='+tk+' type='+msg.type);

    // Phase detection
    if (msg.type==='system'||msg.type==='phase') {
        const c = (msg.content||'').toLowerCase();
        ['research','planning','discussion','development','template_selection','qa'].forEach(p => {
            if (c.includes(p) && c.includes('phase')) { L.phase('Detected phase change: '+p); setPhase(p); updateWB('Phase: '+p.toUpperCase()); }
        });
    }

    if (fk && AGENTS[fk]) {
        if (!S.agents[fk]) S.agents[fk] = {state:'idle',progress:0,tasks:[],messages:[]};
        const prevState = S.agents[fk].state;
        S.agents[fk].state = 'working';
        if (prevState !== 'working') L.agent(fk+' state: '+prevState+' -> working');

        if (msg.type==='response'||msg.type==='delegate') showBubble(fk, summary(msg.content));
        if (msg.type==='delegate' && tk && AGENTS[tk]) { L.agent('Delegation: '+fk+' -> '+tk); sendAirplane(fk,tk,AGENTS[fk].color); }
        if (msg.type==='file_create'||(msg.metadata&&msg.metadata.file_path)) {
            const fname = (msg.metadata?.file_path||'file').split('/').pop();
            L.agent(fk+' created file: '+fname);
            addPostIt(fk, fname);
        }
        if (msg.type==='response' && (fk==='ceo'||fk==='architect')) {
            const s = summary(msg.content);
            if (s.length>5) updateWB(s);
        }
    } else if (msg.from) {
        L.warn('Unknown agent role: from='+msg.from+' (normalized='+fk+')');
    }

    if (S.selected===fk||S.selected===tk) openPanel(S.selected);
    render();
}

function handleTaskEvt(evt) {
    const k = normRole(evt.agent_role);
    if (!k) { L.warn('handleTaskEvt: unknown role='+evt.agent_role); return; }
    if (!S.agents[k]) S.agents[k] = {state:'idle',progress:0,tasks:[],messages:[]};
    const a = S.agents[k];
    a.tasks = a.tasks || [];

    switch(evt.event_type) {
        case 'task_started':
            L.task(k+' STARTED: '+(evt.data?.title||'untitled')+' id='+evt.task_id+' type='+(evt.data?.task_type||''));
            a.state='working'; a.taskTitle=evt.data?.title||''; a.progress=0;
            a.tasks.push({id:evt.task_id,title:evt.data?.title||'',task_type:evt.data?.task_type||'',status:'in_progress',progress:0});
            showBubble(k,'Starting: '+(evt.data?.title||'task'));
            break;
        case 'task_progress':
            L.task(k+' PROGRESS: '+(evt.data?.progress||0)+'% id='+evt.task_id);
            a.state='working'; a.progress=evt.data?.progress||0;
            const pt = a.tasks.find(t=>t.id===evt.task_id);
            if(pt) pt.progress = evt.data?.progress||0;
            break;
        case 'task_completed':
            L.task(k+' COMPLETED: id='+evt.task_id+' outputs='+(evt.data?.outputs?.length||0));
            a.state='completed'; a.progress=100;
            const ct = a.tasks.find(t=>t.id===evt.task_id);
            if(ct){ct.status='completed';ct.progress=100;}
            showBubble(k,'Done!');
            playSound('complete');
            addXP(k, 25);
            gameStats.files += (evt.data?.outputs?.length || 1);
            saveStats(); updateStats();
            checkAchievement('full_house');
            setTimeout(()=>{ if(S.agents[k]?.state==='completed'){S.agents[k].state='idle';render();} },5000);
            break;
        case 'task_failed':
            L.error(k+' FAILED: id='+evt.task_id+' data='+JSON.stringify(evt.data||{}));
            a.state='error';
            const ft = a.tasks.find(t=>t.id===evt.task_id);
            if(ft) ft.status='failed';
            showBubble(k,'Error!');
            playSound('error');
            break;
        default:
            L.warn('Unknown task event: '+evt.event_type, evt);
    }
    render();
    if(S.selected===k) openPanel(k);
}

// ============================
// API
// ============================
async function submitTask() {
    const inp = document.getElementById('taskInp');
    const task = inp.value.trim();
    if (!task) return;
    const btn = document.getElementById('taskBtn');
    btn.disabled = true; S.running = true;
    L.task('Submitting task: project='+S.projectId+' task="'+task+'"');
    document.getElementById('phaseText').textContent = 'STARTING';
    S.agents = {}; S.wbNotes = []; S.messages = [];
    document.getElementById('msgLog').innerHTML = '';
    document.querySelectorAll('.postit,.bubble').forEach(e=>e.remove());
    updateWB('Task: '+task.substring(0,20));
    render();

    try {
        const res = await fetch('/api/task',{method:'POST',headers:{'Content-Type':'application/json'},
            body:JSON.stringify({task,project_id:S.projectId})});
        const d = await res.json();
        L.api('POST /api/task response:', d);
        if (d.success) {
            L.task('Task accepted: id='+d.task_id+' project='+d.project_id);
            inp.value='';
            buildStartTime = Date.now();
            addLog({from:'system',content:'Task: '+task,timestamp:new Date().toISOString()});
            checkAchievement('night_owl');
        }
        else { L.error('Task rejected:', d.error||'Unknown'); addLog({from:'system',content:'Error: '+(d.error||'Failed'),timestamp:new Date().toISOString()}); btn.disabled=false; S.running=false; }
    } catch(e) { L.error('Task submit failed:', e); addLog({from:'system',content:'Connection error',timestamp:new Date().toISOString()}); btn.disabled=false; S.running=false; }
}

async function fetchProjects() {
    try {
        const d = await (await fetch('/api/projects')).json();
        L.api('GET /api/projects:', (d.projects||[]).length, 'projects');
        const sel = document.getElementById('projectSelect');
        sel.innerHTML = '';
        (d.projects||[]).forEach(p => {
            const o = document.createElement('option');
            o.value=p.id; o.textContent=p.id;
            if(p.id===S.projectId) o.selected=true;
            sel.appendChild(o);
        });
    } catch(e){ L.error('fetchProjects failed:', e); }
}

async function fetchStatus() {
    try {
        const d = await (await fetch('/api/status')).json();
        const was = S.running;
        S.running = d.task_running;
        if (was !== d.task_running) L.info('Status changed: running='+was+' -> '+d.task_running+' projects='+JSON.stringify(d.running_projects)+' ws_clients='+d.ws_clients+' msgs='+d.message_count);
        if (was && !d.task_running) {
            L.task('BUILD COMPLETE! Duration='+(buildStartTime ? ((Date.now()-buildStartTime)/1000).toFixed(1)+'s' : 'unknown'));
            document.getElementById('taskBtn').disabled = false;
            document.getElementById('phaseText').textContent = 'COMPLETE';
            gameStats.builds++;
            saveStats(); updateStats();
            checkAchievement('first_build');
            checkAchievement('builds5');
            if (buildStartTime && (Date.now() - buildStartTime) < 180000) checkAchievement('speed_demon');
            buildStartTime = null;
            Object.keys(S.agents).forEach(k => { if(S.agents[k].state==='working') S.agents[k].state='completed'; });
            render();
            setTimeout(() => {
                if(!S.running) {
                    document.getElementById('phaseText').textContent='IDLE';
                    Object.keys(S.agents).forEach(k => { S.agents[k].state='idle'; });
                    render();
                }
            }, 8000);
        }
    } catch(e){ L.error('fetchStatus failed:', e); }
}

async function fetchAgentTasks() {
    try {
        const d = await (await fetch('/api/agent-tasks?project='+S.projectId)).json();
        if (!d.tasks) { L.debug('fetchAgentTasks: no tasks'); return; }
        if (d.tasks.length > 0) L.api('GET /api/agent-tasks:', d.count, 'tasks');
        d.tasks.forEach(task => {
            const k = normRole(task.agent_role);
            if (!k || !AGENTS[k]) { L.warn('fetchAgentTasks: unknown role='+task.agent_role); return; }
            if (!S.agents[k]) S.agents[k] = {state:'idle',progress:0,tasks:[],messages:[]};
            S.agents[k].tasks = S.agents[k].tasks || [];
            const td = {id:task.id,title:task.title,task_type:task.task_type,status:task.status,progress:task.progress};
            const idx = S.agents[k].tasks.findIndex(t=>t.id===task.id);
            if(idx>=0) S.agents[k].tasks[idx]=td; else S.agents[k].tasks.push(td);
            if(task.status==='in_progress'){S.agents[k].state='working';S.agents[k].taskTitle=task.title;S.agents[k].progress=task.progress;}
            else if(task.status==='completed'&&S.agents[k].state!=='working') S.agents[k].state='completed';
            else if(task.status==='failed') { L.error('Agent task failed: '+k+' id='+task.id+' title='+task.title); S.agents[k].state='error'; }
        });
        render();
    } catch(e){ L.error('fetchAgentTasks failed:', e); }
}

// ============================
// UTILS
// ============================
function normRole(r) {
    if(!r) return null;
    const m = {ceo:'ceo',pm:'pm',ux:'ux',ui:'ui',security:'security',architect:'architect',
        'senior-dev':'senior-dev','senior_dev':'senior-dev',seniordev:'senior-dev',
        'junior-dev':'junior-dev','junior_dev':'junior-dev',juniordev:'junior-dev'};
    return m[r.toLowerCase().replace(/_/g,'-')] || null;
}
function summary(c) {
    if(!c) return '';
    const s = c.replace(/[#*]/g,'').trim().split(/[.\n]/)[0].trim();
    return s.length>50?s.substring(0,47)+'...':s;
}
function esc(s) { return s?s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'):''; }
function fmtTime(t) {
    if(!t) return '';
    return new Date(t).toLocaleTimeString([],{hour:'2-digit',minute:'2-digit',second:'2-digit'});
}

// Clock
function tickClock() {
    const now = new Date();
    const h = (now.getHours() % 12) * 30 + now.getMinutes() * 0.5;
    const m = now.getMinutes() * 6;
    const hEl = document.getElementById('clockH');
    const mEl = document.getElementById('clockM');
    if (hEl) hEl.style.transform = 'translateX(-50%) rotate(' + h + 'deg)';
    if (mEl) mEl.style.transform = 'translateX(-50%) rotate(' + m + 'deg)';
}

// Phase banner
function showPhaseBanner(text, color) {
    const banner = document.getElementById('phaseBanner');
    if (!banner) return;
    banner.textContent = text;
    banner.style.background = color || 'rgba(78,205,196,0.9)';
    banner.style.color = '#fff';
    banner.classList.add('show');
    setTimeout(() => banner.classList.remove('show'), 2500);
}

// Random blink
function randomBlink() {
    document.querySelectorAll('.face').forEach(f => {
        if(Math.random()>.7){f.classList.add('blink');setTimeout(()=>f.classList.remove('blink'),200);}
    });
}

// Per-agent idle behaviors
const IDLE_BEHAVIORS = {
    ceo: ['idle-lean', 'idle-nod'],
    pm: ['idle-tap', 'idle-look-right'],
    ux: ['idle-look-left', 'idle-look-right', 'idle-nod'],
    ui: ['idle-look-left', 'idle-stretch'],
    security: ['idle-scan', 'idle-look-left', 'idle-look-right'],
    architect: ['idle-nod', 'idle-lean'],
    'senior-dev': ['idle-lean', 'idle-look-right'],
    'junior-dev': ['idle-peek', 'idle-look-left', 'idle-nod']
};

function randomIdleBehavior() {
    document.querySelectorAll('.agent').forEach(agentEl => {
        const char = agentEl.querySelector('.char');
        const key = agentEl.dataset.agent;
        if (!char || char.classList.contains('typing') || Math.random() > 0.35) return;

        const behaviors = IDLE_BEHAVIORS[key] || ['idle-look-left'];
        const pick = behaviors[Math.floor(Math.random() * behaviors.length)];

        // For peek, move the whole agent element
        if (pick === 'idle-peek') {
            agentEl.classList.add('idle-peek');
            setTimeout(() => agentEl.classList.remove('idle-peek'), 2000);
        } else {
            char.classList.add(pick);
            setTimeout(() => char.classList.remove(pick), 1800);
        }
    });
}

// Cat - click to wake, roam, sleep
let catRoaming = false;
function initCatClick() {
    const cat = document.getElementById('officeCat');
    if (!cat) return;
    cat.addEventListener('click', () => {
        if (catRoaming) return;
        catRoaming = true;
        catClicks++;
        localStorage.setItem('ah_cat_clicks', catClicks);
        checkAchievement('cat_lover');
        playSound('meow');
        const zzz = cat.querySelector('.cat-zzz');
        if (zzz) zzz.textContent = 'Meow!';
        cat.classList.add('awake');

        // Stretch
        setTimeout(() => {
            cat.classList.add('roaming');
            // Pick a random agent desk to visit
            const keys = Object.keys(AGENTS);
            const target = AGENTS[keys[Math.floor(Math.random() * keys.length)]];
            cat.style.left = (target.deskX + 60) + 'px';
            cat.style.top = (target.deskY - 130) + 'px'; // relative to floor

            // Arrive, sit, then go back to sleep
            setTimeout(() => {
                if (zzz) zzz.textContent = '!';
                setTimeout(() => {
                    // Walk back home
                    cat.style.left = '560px';
                    cat.style.top = '50px';
                    setTimeout(() => {
                        cat.classList.remove('awake', 'roaming');
                        if (zzz) zzz.textContent = 'z';
                        catRoaming = false;
                    }, 3000);
                }, 2000);
            }, 3500);
        }, 500);
    });
}

// Coffee machine
let coffeeCount = parseInt(localStorage.getItem('ah_coffees') || '0');
function initCoffeeClick() {
    const cm = document.getElementById('coffeeMachine');
    if (!cm) return;
    updateBrewCount();
    cm.addEventListener('click', () => {
        if (cm.classList.contains('brewing')) return;
        cm.classList.add('brewing');
        playSound('brew');
        coffeeCount++;
        localStorage.setItem('ah_coffees', coffeeCount);
        updateBrewCount();
        checkAchievement('coffee10');
        setTimeout(() => {
            cm.classList.remove('brewing');
            // Random agent gets a bubble
            const keys = Object.keys(AGENTS);
            const k = keys[Math.floor(Math.random() * keys.length)];
            showBubble(k, 'Coffee time!');
        }, 2500);
    });
}
function updateBrewCount() {
    const el = document.getElementById('brewCount');
    if (el) el.textContent = coffeeCount > 0 ? coffeeCount + ' brewed' : 'COFFEE';
}

// Whiteboard modal
function initWhiteboardClick() {
    document.getElementById('whiteboard').addEventListener('click', () => {
        const modal = document.createElement('div');
        modal.className = 'wb-modal';
        const entries = S.wbNotes.length ? S.wbNotes.map(n =>
            '<div class="wb-entry">'+esc(n)+'</div>'
        ).join('') : '<div style="color:#999;font-style:italic">No notes yet. Start a task!</div>';

        modal.innerHTML =
            '<div class="wb-modal-content">' +
            '<div class="wb-modal-header">SPRINT BOARD<button class="wb-modal-close">x</button></div>' +
            '<div class="wb-modal-body">' + entries + '</div></div>';

        document.body.appendChild(modal);
        modal.querySelector('.wb-modal-close').addEventListener('click', () => modal.remove());
        modal.addEventListener('click', e => { if (e.target === modal) modal.remove(); });
    });
}

// Day/night cycle
function updateDayNight() {
    const hour = new Date().getHours();
    let timeClass;
    if (hour >= 6 && hour < 9) timeClass = 'dawn';
    else if (hour >= 9 && hour < 17) timeClass = 'day';
    else if (hour >= 17 && hour < 20) timeClass = 'dusk';
    else timeClass = 'night';

    document.querySelectorAll('.window').forEach(w => {
        w.classList.remove('dawn', 'day', 'dusk', 'night');
        w.classList.add(timeClass);
    });
}

// 8-bit Sound via Web Audio API
let audioCtx;
function getAudioCtx() {
    if (!audioCtx) {
        try { audioCtx = new (window.AudioContext || window.webkitAudioContext)(); } catch(e) {}
    }
    return audioCtx;
}

function playSound(type) {
    const ctx = getAudioCtx();
    if (!ctx) return;
    try {
        const osc = ctx.createOscillator();
        const gain = ctx.createGain();
        osc.connect(gain);
        gain.connect(ctx.destination);
        gain.gain.value = 0.06; // quiet

        switch(type) {
            case 'blip':
                osc.type = 'square';
                osc.frequency.setValueAtTime(800, ctx.currentTime);
                osc.frequency.exponentialRampToValueAtTime(400, ctx.currentTime + 0.08);
                gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.1);
                osc.start(ctx.currentTime);
                osc.stop(ctx.currentTime + 0.1);
                break;
            case 'whoosh':
                osc.type = 'sawtooth';
                osc.frequency.setValueAtTime(200, ctx.currentTime);
                osc.frequency.exponentialRampToValueAtTime(800, ctx.currentTime + 0.15);
                gain.gain.value = 0.03;
                gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.2);
                osc.start(ctx.currentTime);
                osc.stop(ctx.currentTime + 0.2);
                break;
            case 'meow':
                osc.type = 'sine';
                osc.frequency.setValueAtTime(600, ctx.currentTime);
                osc.frequency.exponentialRampToValueAtTime(900, ctx.currentTime + 0.1);
                osc.frequency.exponentialRampToValueAtTime(500, ctx.currentTime + 0.3);
                gain.gain.value = 0.08;
                gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.35);
                osc.start(ctx.currentTime);
                osc.stop(ctx.currentTime + 0.35);
                break;
            case 'brew':
                osc.type = 'triangle';
                osc.frequency.setValueAtTime(150, ctx.currentTime);
                gain.gain.value = 0.04;
                gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.8);
                osc.start(ctx.currentTime);
                osc.stop(ctx.currentTime + 0.8);
                break;
            case 'phase':
                // 4-note arpeggio
                [440, 554, 659, 880].forEach((freq, i) => {
                    const o = ctx.createOscillator();
                    const g = ctx.createGain();
                    o.type = 'square';
                    o.frequency.value = freq;
                    o.connect(g);
                    g.connect(ctx.destination);
                    g.gain.value = 0.05;
                    g.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.15*(i+1) + 0.1);
                    o.start(ctx.currentTime + 0.15*i);
                    o.stop(ctx.currentTime + 0.15*(i+1) + 0.1);
                });
                return;
            case 'complete':
                [523, 659, 784, 1047].forEach((freq, i) => {
                    const o = ctx.createOscillator();
                    const g = ctx.createGain();
                    o.type = 'square';
                    o.frequency.value = freq;
                    o.connect(g);
                    g.connect(ctx.destination);
                    g.gain.value = 0.06;
                    g.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.2*(i+1) + 0.15);
                    o.start(ctx.currentTime + 0.2*i);
                    o.stop(ctx.currentTime + 0.2*(i+1) + 0.15);
                });
                return;
            case 'error':
                osc.type = 'square';
                osc.frequency.setValueAtTime(200, ctx.currentTime);
                osc.frequency.setValueAtTime(150, ctx.currentTime + 0.15);
                gain.gain.value = 0.06;
                gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.3);
                osc.start(ctx.currentTime);
                osc.stop(ctx.currentTime + 0.3);
                break;
        }
    } catch(e) {}
}

// ============================
// XP & LEVEL SYSTEM
// ============================
const XP_PER_LEVEL = 100;
let agentXP = JSON.parse(localStorage.getItem('ah_xp') || '{}');

function addXP(agentKey, amount) {
    if (!agentXP[agentKey]) agentXP[agentKey] = { xp: 0, level: 1 };
    const prev = agentXP[agentKey].level;
    agentXP[agentKey].xp += amount;
    const newLevel = Math.floor(agentXP[agentKey].xp / XP_PER_LEVEL) + 1;
    agentXP[agentKey].level = newLevel;
    L.xp(agentKey+' +'+amount+'XP (total='+agentXP[agentKey].xp+' level='+newLevel+')');
    localStorage.setItem('ah_xp', JSON.stringify(agentXP));

    if (newLevel > prev) {
        L.xp(agentKey+' LEVEL UP! '+prev+' -> '+newLevel);
        showLevelUp(agentKey, newLevel);
        checkAchievement('level5', agentKey);
    }
}

function getXP(agentKey) {
    return agentXP[agentKey] || { xp: 0, level: 1 };
}

function showLevelUp(agentKey, level) {
    const ag = AGENTS[agentKey];
    if (!ag) return;
    const scene = document.getElementById('scene');
    const el = document.createElement('div');
    el.className = 'level-up-pop';
    el.style.cssText = 'left:'+(ag.deskX+50)+'px;top:'+(ag.deskY-10)+'px;';
    el.textContent = 'LVL '+level+'!';
    scene.appendChild(el);
    playSound('complete');
    setTimeout(() => el.remove(), 2200);
}

// ============================
// ACHIEVEMENT SYSTEM
// ============================
const ACHIEVEMENTS = {
    first_build:  { icon: '&#9733;', name: 'First Build', desc: 'Complete your first project' },
    coffee10:     { icon: '&#9749;', name: 'Coffee Addict', desc: 'Brew 10 coffees' },
    night_owl:    { icon: '&#9790;', name: 'Night Owl', desc: 'Use the dashboard after midnight' },
    full_house:   { icon: '&#9819;', name: 'Full House', desc: 'All 8 agents working at once' },
    cat_lover:    { icon: '&#9829;', name: 'Cat Person', desc: 'Pet Mr. Whiskers 10 times' },
    level5:       { icon: '&#9889;', name: 'Power Up', desc: 'Get any agent to level 5' },
    speed_demon:  { icon: '&#9889;', name: 'Speed Demon', desc: 'Complete a build in under 3 min' },
    crunch_mode:  { icon: '&#9760;', name: 'Crunch Time', desc: 'Activate the secret code' },
    builds5:      { icon: '&#9734;', name: 'Veteran', desc: 'Complete 5 builds' }
};

let unlockedAchievements = JSON.parse(localStorage.getItem('ah_achievements') || '[]');
let catClicks = parseInt(localStorage.getItem('ah_cat_clicks') || '0');

function checkAchievement(id, extra) {
    if (unlockedAchievements.includes(id)) return;

    let earned = false;
    switch(id) {
        case 'first_build': earned = true; break;
        case 'coffee10': earned = coffeeCount >= 10; break;
        case 'night_owl':
            const h = new Date().getHours();
            earned = h >= 0 && h < 5;
            break;
        case 'full_house':
            const working = Object.values(S.agents).filter(a => a.state === 'working').length;
            earned = working >= 8;
            break;
        case 'cat_lover': earned = catClicks >= 10; break;
        case 'level5':
            earned = extra && agentXP[extra] && agentXP[extra].level >= 5;
            break;
        case 'crunch_mode': earned = true; break;
        case 'builds5': earned = gameStats.builds >= 5; break;
    }

    if (earned) {
        unlockedAchievements.push(id);
        localStorage.setItem('ah_achievements', JSON.stringify(unlockedAchievements));
        showAchievementToast(id);
        playSound('complete');
    }
}

function showAchievementToast(id) {
    const ach = ACHIEVEMENTS[id];
    if (!ach) return;
    const el = document.createElement('div');
    el.className = 'achievement-toast';
    el.innerHTML =
        '<div class="achievement-icon">'+ach.icon+'</div>' +
        '<div class="achievement-info">' +
        '<div class="achievement-label">Achievement Unlocked</div>' +
        '<div class="achievement-name">'+ach.name+'</div>' +
        '<div class="achievement-desc">'+ach.desc+'</div></div>';
    document.body.appendChild(el);
    setTimeout(() => el.remove(), 4200);
}

// ============================
// GAME STATS
// ============================
let gameStats = JSON.parse(localStorage.getItem('ah_stats') || '{"builds":0,"files":0,"messages":0,"uptime":0}');
let buildStartTime = null;

function updateStats() {
    document.getElementById('statBuilds').textContent = gameStats.builds;
    document.getElementById('statFiles').textContent = gameStats.files;
    const maxLvl = Math.max(1, ...Object.values(agentXP).map(x => x.level || 1));
    document.getElementById('statLevel').textContent = maxLvl;
}

function saveStats() {
    localStorage.setItem('ah_stats', JSON.stringify(gameStats));
}

// ============================
// KONAMI CODE
// ============================
let konamiBuffer = [];
const KONAMI = [38,38,40,40,37,39,37,39,66,65]; // up up down down left right left right B A
let crunchActive = false;

function initKonami() {
    document.addEventListener('keydown', e => {
        konamiBuffer.push(e.keyCode);
        if (konamiBuffer.length > 10) konamiBuffer.shift();
        if (konamiBuffer.length === 10 && konamiBuffer.every((v,i) => v === KONAMI[i])) {
            activateCrunchMode();
            konamiBuffer = [];
        }
    });
}

function activateCrunchMode() {
    if (crunchActive) return;
    crunchActive = true;
    checkAchievement('crunch_mode');

    // Banner flash
    const banner = document.createElement('div');
    banner.className = 'crunch-banner';
    banner.textContent = 'CRUNCH MODE';
    document.body.appendChild(banner);
    setTimeout(() => banner.remove(), 1800);

    // Apply crunch styles
    document.body.classList.add('crunch-mode');

    // Speed up all typing agents
    document.querySelectorAll('.typing .arm').forEach(arm => {
        arm.style.animationDuration = '0.15s';
    });

    playSound('error');
    setTimeout(() => playSound('phase'), 300);

    // Deactivate after 15 seconds
    setTimeout(() => {
        document.body.classList.remove('crunch-mode');
        crunchActive = false;
        document.querySelectorAll('.arm').forEach(arm => {
            arm.style.animationDuration = '';
        });
    }, 15000);
}

// ============================
// NEW PROJECT MODAL
// ============================
function openNewProjectModal() {
    const modal = document.createElement('div');
    modal.className = 'np-modal';
    modal.innerHTML =
        '<div class="np-box">' +
        '<div class="np-title">NEW PROJECT</div>' +
        '<input type="text" class="np-input" id="npInput" placeholder="my-awesome-app" autofocus>' +
        '<div class="np-hint">LOWERCASE, HYPHENS OK, NO SPACES</div>' +
        '<div class="np-error" id="npError"></div>' +
        '<div class="np-actions">' +
        '<button class="np-btn cancel" id="npCancel">CANCEL</button>' +
        '<button class="np-btn create" id="npCreate">CREATE</button>' +
        '</div></div>';
    document.body.appendChild(modal);

    const inp = document.getElementById('npInput');
    const errEl = document.getElementById('npError');

    function showErr(msg) { errEl.textContent = msg; errEl.style.display = 'block'; }
    function hideErr() { errEl.style.display = 'none'; }

    function doCreate() {
        const raw = inp.value.trim();
        // Sanitize: lowercase, replace spaces with hyphens, remove invalid chars
        const name = raw.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9\-_]/g, '');
        if (!name) { showErr('PROJECT NAME REQUIRED'); return; }
        if (name.length > 40) { showErr('NAME TOO LONG (MAX 40)'); return; }

        L.info('Creating new project:', name);

        // Add to select and switch to it
        const sel = document.getElementById('projectSelect');
        // Check if already exists
        const exists = Array.from(sel.options).some(o => o.value === name);
        if (!exists) {
            const opt = document.createElement('option');
            opt.value = name; opt.textContent = name;
            sel.appendChild(opt);
        }
        sel.value = name;
        S.projectId = name;
        S.agents = {}; S.messages = [];
        document.getElementById('msgLog').innerHTML = '';
        document.querySelectorAll('.postit,.bubble').forEach(el => el.remove());
        render();
        fetchAgentTasks();
        modal.remove();
        L.info('Switched to project:', name);
    }

    document.getElementById('npCancel').addEventListener('click', () => modal.remove());
    document.getElementById('npCreate').addEventListener('click', doCreate);
    inp.addEventListener('keydown', e => {
        hideErr();
        if (e.key === 'Enter') doCreate();
        if (e.key === 'Escape') modal.remove();
    });
    modal.addEventListener('click', e => { if (e.target === modal) modal.remove(); });
    inp.focus();
}

// ============================
// INIT
// ============================
function init() {
    L.info('=== 8-Bit Office Init ===');
    L.info('Project:', S.projectId);
    L.info('Stored stats:', JSON.stringify(gameStats));
    L.info('Stored XP:', JSON.stringify(agentXP));
    L.info('Achievements:', JSON.stringify(unlockedAchievements));
    render();
    connectWS();
    fetchProjects();
    tickClock();
    updateDayNight();
    updateStats();
    initCatClick();
    initCoffeeClick();
    initWhiteboardClick();
    initKonami();
    L.info('Init complete, polling: status=3s agentTasks=5s');

    document.getElementById('taskBtn').addEventListener('click', submitTask);
    document.getElementById('taskInp').addEventListener('keydown', e => { if(e.key==='Enter') submitTask(); });
    document.getElementById('spClose').addEventListener('click', closePanel);
    document.getElementById('newProjBtn').addEventListener('click', openNewProjectModal);
    document.getElementById('projectSelect').addEventListener('change', e => {
        S.projectId = e.target.value;
        S.agents = {}; S.messages = [];
        document.getElementById('msgLog').innerHTML = '';
        document.querySelectorAll('.postit,.bubble').forEach(el=>el.remove());
        render(); fetchAgentTasks();
    });
    document.addEventListener('keydown', e => { if(e.key==='Escape') closePanel(); });

    // Click outside panel to close
    document.querySelector('.office-wrap').addEventListener('click', e => {
        if (!e.target.closest('.desk-unit') && !e.target.closest('.agent')) closePanel();
    });

    setInterval(tickClock, 1000);
    setInterval(updateDayNight, 60000);
    setInterval(fetchStatus, 3000);
    setInterval(fetchAgentTasks, 5000);
    setInterval(randomBlink, 2500);
    setInterval(randomIdleBehavior, 4000);

    fetchAgentTasks();
}

document.addEventListener('DOMContentLoaded', init);
</script>
</body>
</html>` + "\n"
