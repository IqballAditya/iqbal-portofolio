// ==============================================================================
// ADITH AI STUDIO - Figma UI Generator Plugin
// ==============================================================================

(async () => {
    console.log("🚀 Generating ADITH AI STUDIO UI canvas in Figma...");

    const rgb = (r, g, b) => ({ r: r / 255, g: g / 255, b: b / 255 });
    const solid = (r, g, b, opacity = 1) => [{ type: 'SOLID', color: rgb(r, g, b), opacity }];

    const T = {
        bgApp: rgb(13, 15, 18),        // #0D0F12 Deep dark canvas
        bgPanel: rgb(22, 25, 32),      // #161920 Card / Panel background
        bgPanelLight: rgb(30, 35, 45), // #1E232D Lighter container
        bgInput: rgb(17, 20, 26),      // #11141A Input background
        border: rgb(40, 46, 58),       // #282E3A Border lines
        borderLight: rgb(55, 63, 80),  // #373F50 Highlight border
        textWhite: rgb(255, 255, 255),
        textMuted: rgb(150, 160, 178), // Subtitles / dim text
        textDim: rgb(100, 110, 130),
        accentIndigo: rgb(99, 102, 241),// #6366F1 Brand Indigo
        accentEmerald: rgb(16, 185, 129),// #10B981 Success / Diff Add
        accentAmber: rgb(245, 158, 11),  // #F59E0B Warning / Tokens
        accentRose: rgb(244, 63, 94),    // #F43F5E Diff Remove
    };

    try {
        await figma.loadFontAsync({ family: "Inter", style: "Regular" });
        await figma.loadFontAsync({ family: "Inter", style: "Medium" });
        await figma.loadFontAsync({ family: "Inter", style: "Semi Bold" });
        await figma.loadFontAsync({ family: "Inter", style: "Bold" });
    } catch (e) {
        console.warn("Font loading fallback: ", e);
    }

    const createTxt = (chars, size = 14, weight = "Regular", color = T.textWhite) => {
        const t = figma.createText();
        try {
            t.fontName = { family: "Inter", style: weight };
        } catch (_) {}
        t.characters = chars;
        t.fontSize = size;
        t.fills = [{ type: 'SOLID', color }];
        return t;
    };

    const createAL = (mode, gap = 0, padH = 0, padV = 0, fills = null) => {
        const f = figma.createFrame();
        f.layoutMode = mode;
        f.itemSpacing = gap;
        f.paddingLeft = padH;
        f.paddingRight = padH;
        f.paddingTop = padV;
        f.paddingBottom = padV;
        f.primaryAxisSizingMode = "AUTO";
        f.counterAxisSizingMode = "AUTO";
        if (fills) f.fills = fills;
        else f.fills = [];
        return f;
    };

    let startX = 0;
    for (const child of figma.currentPage.children) {
        if (child.x + child.width > startX) {
            startX = child.x + child.width + 120;
        }
    }

    // 1. ROOT STUDIO CANVAS
    const root = figma.createFrame();
    root.name = "ADITH AI STUDIO - Web AI IDE";
    root.resize(1440, 920);
    root.x = startX;
    root.y = 0;
    root.fills = solid(13, 15, 18);
    root.layoutMode = "VERTICAL";
    root.primaryAxisSizingMode = "FIXED";
    root.counterAxisSizingMode = "FIXED";
    root.itemSpacing = 0;
    figma.currentPage.appendChild(root);

    // 2. HEADER TOPBAR
    const header = createAL("HORIZONTAL", 16, 24, 14, solid(22, 25, 32));
    header.resize(1440, 64);
    header.layoutSizingHorizontal = "FIXED";
    header.counterAxisAlignItems = "CENTER";
    header.strokes = solid(40, 46, 58);
    header.strokeWeight = 1;

    const brandGrp = createAL("HORIZONTAL", 10, 0, 0);
    brandGrp.counterAxisAlignItems = "CENTER";
    const brandBadge = createAL("HORIZONTAL", 0, 8, 4, solid(99, 102, 241));
    brandBadge.cornerRadius = 6;
    brandBadge.appendChild(createTxt("⚡ ADITH", 12, "Bold", T.textWhite));
    brandGrp.appendChild(brandBadge);
    brandGrp.appendChild(createTxt("AI STUDIO", 15, "Bold", T.textWhite));
    brandGrp.appendChild(createTxt("v1.0.0 • Web Cloud IDE", 11, "Regular", T.textDim));
    header.appendChild(brandGrp);

    const midStatus = createAL("HORIZONTAL", 8, 12, 6, solid(30, 35, 45));
    midStatus.cornerRadius = 20;
    midStatus.strokes = solid(55, 63, 80);
    midStatus.strokeWeight = 1;
    midStatus.appendChild(createTxt("● Claude 3.5 Sonnet + R1 (Agent Active)", 12, "Medium", T.accentEmerald));
    header.appendChild(midStatus);

    const rightActions = createAL("HORIZONTAL", 12, 0, 0);
    rightActions.counterAxisAlignItems = "CENTER";
    const runBtn = createAL("HORIZONTAL", 6, 14, 8, solid(99, 102, 241));
    runBtn.cornerRadius = 6;
    runBtn.appendChild(createTxt("▶ Run Agent (Ctrl+Enter)", 12, "Semi Bold", T.textWhite));
    rightActions.appendChild(runBtn);
    header.appendChild(rightActions);
    root.appendChild(header);

    // 3. WORKSPACE
    const workspace = createAL("HORIZONTAL", 0, 0, 0);
    workspace.resize(1440, 856);
    workspace.layoutSizingHorizontal = "FIXED";
    workspace.layoutSizingVertical = "FIXED";

    // 4. CHAT PANEL (LEFT)
    const chatPanel = createAL("VERTICAL", 12, 16, 16, solid(22, 25, 32));
    chatPanel.resize(440, 856);
    chatPanel.layoutSizingHorizontal = "FIXED";
    chatPanel.layoutSizingVertical = "FIXED";
    chatPanel.strokes = solid(40, 46, 58);
    chatPanel.strokeWeight = 1;

    const chatHeader = createAL("HORIZONTAL", 8, 0, 0);
    chatHeader.counterAxisAlignItems = "CENTER";
    chatHeader.appendChild(createTxt("💬 ADITH AI AGENT COPILOT", 13, "Bold", T.textWhite));
    chatPanel.appendChild(chatHeader);

    const promptCard = createAL("VERTICAL", 6, 12, 10, solid(30, 35, 45));
    promptCard.cornerRadius = 8;
    promptCard.layoutSizingHorizontal = "FILL";
    promptCard.appendChild(createTxt("USER (Iqbal Aditya)", 11, "Semi Bold", T.accentIndigo));
    promptCard.appendChild(createTxt("Optimalkan worker pool Golang dan hubungkan AST parser dengan Qdrant vector memory.", 12, "Regular", T.textWhite));
    chatPanel.appendChild(promptCard);

    const agentCard = createAL("VERTICAL", 8, 12, 10, solid(17, 20, 26));
    agentCard.cornerRadius = 8;
    agentCard.layoutSizingHorizontal = "FILL";
    agentCard.strokes = solid(99, 102, 241, 0.4);
    agentCard.strokeWeight = 1;
    agentCard.appendChild(createTxt("🧠 REASONING TRACE (Tree-sitter + Qdrant)", 11, "Semi Bold", T.accentEmerald));
    agentCard.appendChild(createTxt("1. Membaca AST internal/worker/pool.go...\n2. Menemukan lock contention pada channel dispatch.\n3. Mengganti mutex global dengan lock-free ring buffer.\n4. Mengindeks chunk baru ke vector collection 'codebase_main'.", 11, "Regular", T.textMuted));
    chatPanel.appendChild(agentCard);

    const diffCard = createAL("VERTICAL", 6, 12, 10, solid(13, 15, 18));
    diffCard.cornerRadius = 8;
    diffCard.layoutSizingHorizontal = "FILL";
    diffCard.appendChild(createTxt("PROPOSED DIFF (pool.go)", 11, "Semi Bold", T.textDim));
    diffCard.appendChild(createTxt("- func (p *Pool) Dispatch(j Job) { p.mu.Lock(); ... }", 11, "Regular", T.accentRose));
    diffCard.appendChild(createTxt("+ func (p *Pool) Dispatch(j Job) bool { return p.ring.Push(j) }", 11, "Semi Bold", T.accentEmerald));
    chatPanel.appendChild(diffCard);

    const chatInput = createAL("HORIZONTAL", 8, 12, 12, solid(17, 20, 26));
    chatInput.cornerRadius = 8;
    chatInput.layoutSizingHorizontal = "FILL";
    chatInput.strokes = solid(55, 63, 80);
    chatInput.strokeWeight = 1;
    chatInput.appendChild(createTxt("Instruksikan ADITH AI... (@file, #doc, /tool)", 12, "Regular", T.textDim));
    chatPanel.appendChild(chatInput);
    workspace.appendChild(chatPanel);

    // 5. CODE & TERMINAL PANE (RIGHT)
    const codePane = createAL("VERTICAL", 0, 0, 0, solid(13, 15, 18));
    codePane.resize(1000, 856);
    codePane.layoutSizingHorizontal = "FIXED";
    codePane.layoutSizingVertical = "FIXED";

    const editorTabs = createAL("HORIZONTAL", 4, 16, 8, solid(22, 25, 32));
    editorTabs.resize(1000, 42);
    editorTabs.layoutSizingHorizontal = "FIXED";
    const tab1 = createAL("HORIZONTAL", 6, 12, 6, solid(13, 15, 18));
    tab1.cornerRadius = 4;
    tab1.appendChild(createTxt("⚡ pool.go (Modified)", 12, "Medium", T.accentIndigo));
    editorTabs.appendChild(tab1);
    const tab2 = createAL("HORIZONTAL", 6, 12, 6, solid(22, 25, 32));
    tab2.appendChild(createTxt("orchestrator.py", 12, "Regular", T.textDim));
    editorTabs.appendChild(tab2);
    codePane.appendChild(editorTabs);

    const editorArea = createAL("VERTICAL", 6, 24, 16, solid(13, 15, 18));
    editorArea.resize(1000, 564);
    editorArea.layoutSizingHorizontal = "FIXED";
    editorArea.layoutSizingVertical = "FIXED";
    editorArea.appendChild(createTxt("package worker", 12, "Regular", T.accentIndigo));
    editorArea.appendChild(createTxt("// ADITH AI: Lock-free high-throughput ring buffer", 12, "Regular", T.textDim));
    editorArea.appendChild(createTxt("type LockFreePool struct {", 12, "Bold", T.textWhite));
    editorArea.appendChild(createTxt("    ring      *ringbuffer.Buffer", 12, "Regular", T.textMuted));
    editorArea.appendChild(createTxt("    workers   []*WorkerThread", 12, "Regular", T.textMuted));
    editorArea.appendChild(createTxt("    metrics   *atomic.Uint64", 12, "Regular", T.textMuted));
    editorArea.appendChild(createTxt("}", 12, "Bold", T.textWhite));
    codePane.appendChild(editorArea);

    const terminalDrawer = createAL("VERTICAL", 8, 16, 12, solid(22, 25, 32));
    terminalDrawer.resize(1000, 250);
    terminalDrawer.layoutSizingHorizontal = "FIXED";
    terminalDrawer.layoutSizingVertical = "FIXED";
    terminalDrawer.strokes = solid(40, 46, 58);
    terminalDrawer.strokeWeight = 1;
    terminalDrawer.appendChild(createTxt("⚙️ TERMINAL SANDBOX (Firecracker MicroVM)", 11, "Bold", T.textWhite));

    const termBox = createAL("VERTICAL", 4, 12, 10, solid(13, 15, 18));
    termBox.cornerRadius = 6;
    termBox.layoutSizingHorizontal = "FILL";
    termBox.appendChild(createTxt("$ go test -v -bench=. ./internal/worker/...", 11, "Medium", T.textMuted));
    termBox.appendChild(createTxt("--- PASS: TestWorkerPool_Throughput (0.04s) - 5,280 ops/sec sustained", 11, "Semi Bold", T.accentEmerald));
    termBox.appendChild(createTxt("BenchmarkWorkerPool-16    1000000    1840 ns/op    120 B/op", 11, "Regular", T.accentAmber));
    termBox.appendChild(createTxt("PASS: ok  iqbal-portfolio/internal/worker  0.182s", 11, "Bold", T.accentEmerald));
    terminalDrawer.appendChild(termBox);
    codePane.appendChild(terminalDrawer);

    workspace.appendChild(codePane);
    root.appendChild(workspace);

    figma.viewport.scrollAndZoomIntoView([root]);
    figma.notify("✅ ADITH AI STUDIO Frame Generated Successfully!");
    figma.closePlugin();
})();
