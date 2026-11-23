// ==============================================================================
// ADITH AI STUDIO - Figma UI Generator Script
// How to run:
// 1. Open your Figma design file in Google Chrome / Browser.
// 2. Press F12 to open Developer Tools, then click the "Console" tab.
// 3. Paste this entire script into the Console and press Enter.
// ==============================================================================

(async () => {
    console.log("🚀 Generating ADITH AI STUDIO UI canvas in Figma...");

    // Helper: color converter 0-255 to 0-1
    const rgb = (r, g, b) => ({ r: r / 255, g: g / 255, b: b / 255 });
    const solid = (r, g, b, opacity = 1) => [{ type: 'SOLID', color: rgb(r, g, b), opacity }];

    // Design Tokens (Modern AI Dark Studio)
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

    // Preload Inter fonts
    try {
        await figma.loadFontAsync({ family: "Inter", style: "Regular" });
        await figma.loadFontAsync({ family: "Inter", style: "Medium" });
        await figma.loadFontAsync({ family: "Inter", style: "Semi Bold" });
        await figma.loadFontAsync({ family: "Inter", style: "Bold" });
    } catch (e) {
        console.warn("Font loading note: ", e);
    }

    // Helper: Create Text
    const createTxt = (chars, size = 14, weight = "Regular", color = T.textWhite) => {
        const t = figma.createText();
        t.fontName = { family: "Inter", style: weight };
        t.characters = chars;
        t.fontSize = size;
        t.fills = [{ type: 'SOLID', color }];
        return t;
    };

    // Helper: Create Frame with AutoLayout
    const createAL = (mode = "HORIZONTAL", spacing = 8, padH = 0, padV = 0, fills = null) => {
        const f = figma.createFrame();
        f.layoutMode = mode;
        f.itemSpacing = spacing;
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

    // Find clear canvas position
    let startX = 0;
    for (const child of figma.currentPage.children) {
        if (child.x + child.width > startX) {
            startX = child.x + child.width + 120;
        }
    }

    // ==========================================
    // 1. ROOT STUDIO CANVAS (1440 x 920)
    // ==========================================
    const root = figma.createFrame();
    root.name = "ADITH AI STUDIO - Modern Web AI IDE";
    root.resize(1440, 920);
    root.x = startX;
    root.y = 0;
    root.fills = solid(13, 15, 18);
    root.layoutMode = "VERTICAL";
    root.primaryAxisSizingMode = "FIXED";
    root.counterAxisSizingMode = "FIXED";
    root.itemSpacing = 0;
    figma.currentPage.appendChild(root);

    // ==========================================
    // 2. HEADER TOPBAR (1440 x 64)
    // ==========================================
    const header = createAL("HORIZONTAL", 16, 24, 14, solid(22, 25, 32));
    header.resize(1440, 64);
    header.layoutSizingHorizontal = "FIXED";
    header.counterAxisAlignItems = "CENTER";
    header.strokes = solid(40, 46, 58);
    header.strokeWeight = 1;
    header.strokeAlign = "INSIDE";

    // Logo & Brand Pill
    const brandBlock = createAL("HORIZONTAL", 10, 0, 0);
    brandBlock.counterAxisAlignItems = "CENTER";
    const logoBadge = createAL("HORIZONTAL", 0, 8, 4, solid(99, 102, 241));
    logoBadge.cornerRadius = 6;
    logoBadge.appendChild(createTxt("ADITH", 12, "Bold", rgb(255, 255, 255)));
    const brandName = createTxt("AI STUDIO", 15, "Bold", rgb(255, 255, 255));
    brandBlock.appendChild(logoBadge);
    brandBlock.appendChild(brandName);

    // Project & Branch Selector
    const projectPill = createAL("HORIZONTAL", 8, 12, 6, solid(30, 35, 45));
    projectPill.cornerRadius = 20;
    projectPill.strokes = solid(55, 63, 80);
    projectPill.strokeWeight = 1;
    projectPill.appendChild(createTxt("📁 microservice-go", 12, "Medium", T.textWhite));
    projectPill.appendChild(createTxt("🌱 main", 12, "Regular", T.accentEmerald));
    projectPill.appendChild(createTxt("▾", 12, "Regular", T.textMuted));

    // Spacer
    const headerSpacer = figma.createFrame();
    headerSpacer.layoutGrow = 1;
    headerSpacer.fills = [];

    // Model Selector Badge
    const modelPill = createAL("HORIZONTAL", 8, 14, 6, solid(30, 35, 45));
    modelPill.cornerRadius = 20;
    modelPill.strokes = solid(99, 102, 241);
    modelPill.strokeWeight = 1;
    modelPill.appendChild(createTxt("🧠 Claude 3.7 Sonnet (Thinking)", 12, "Semi Bold", rgb(199, 210, 254)));
    modelPill.appendChild(createTxt("▾", 12, "Regular", T.textMuted));

    // Token telemetry badge
    const tokenBadge = createAL("HORIZONTAL", 6, 10, 6, solid(17, 20, 26));
    tokenBadge.cornerRadius = 6;
    tokenBadge.appendChild(createTxt("⚡ 24ms • 1.4k tokens", 11, "Medium", T.accentAmber));

    // Action buttons
    const btnPreview = createAL("HORIZONTAL", 6, 14, 8, solid(30, 35, 45));
    btnPreview.cornerRadius = 6;
    btnPreview.strokes = solid(55, 63, 80);
    btnPreview.strokeWeight = 1;
    btnPreview.appendChild(createTxt("▶ Run Sandbox", 12, "Semi Bold", T.textWhite));

    const btnDeploy = createAL("HORIZONTAL", 6, 16, 8, solid(99, 102, 241));
    btnDeploy.cornerRadius = 6;
    btnDeploy.appendChild(createTxt("🚀 Deploy Web", 12, "Bold", rgb(255, 255, 255)));

    // User Avatar
    const avatar = createAL("HORIZONTAL", 0, 0, 0, solid(16, 185, 129));
    avatar.resize(32, 32);
    avatar.cornerRadius = 16;
    avatar.primaryAxisAlignItems = "CENTER";
    avatar.counterAxisAlignItems = "CENTER";
    avatar.appendChild(createTxt("IA", 12, "Bold", rgb(255, 255, 255)));

    header.appendChild(brandBlock);
    header.appendChild(projectPill);
    header.appendChild(headerSpacer);
    header.appendChild(modelPill);
    header.appendChild(tokenBadge);
    header.appendChild(btnPreview);
    header.appendChild(btnDeploy);
    header.appendChild(avatar);
    root.appendChild(header);

    // ==========================================
    // 3. MAIN WORKSPACE BODY (Split Screen)
    // ==========================================
    const workspace = createAL("HORIZONTAL", 0, 0, 0);
    workspace.resize(1440, 856);
    workspace.layoutSizingHorizontal = "FIXED";
    workspace.layoutSizingVertical = "FILL";

    // ------------------------------------------
    // LEFT PANE: AI CHAT CANVAS (540px)
    // ------------------------------------------
    const chatPane = createAL("VERTICAL", 16, 20, 20, solid(17, 20, 26));
    chatPane.resize(540, 856);
    chatPane.layoutSizingHorizontal = "FIXED";
    chatPane.layoutSizingVertical = "FILL";
    chatPane.strokes = solid(40, 46, 58);
    chatPane.strokeWeight = 1;
    chatPane.strokeAlign = "INSIDE";

    // Chat Header
    const chatHead = createAL("HORIZONTAL", 8, 0, 0);
    chatHead.counterAxisAlignItems = "CENTER";
    chatHead.appendChild(createTxt("💬 Conversational Agent", 14, "Bold", T.textWhite));
    const statusPill = createAL("HORIZONTAL", 0, 8, 2, solid(16, 185, 129, 0.2));
    statusPill.cornerRadius = 10;
    statusPill.appendChild(createTxt("Active Session", 10, "Semi Bold", T.accentEmerald));
    chatHead.appendChild(statusPill);
    chatPane.appendChild(chatHead);

    // Message 1: User Prompt Bubble
    const userMsg = createAL("VERTICAL", 8, 14, 12, solid(30, 35, 45));
    userMsg.cornerRadius = 12;
    userMsg.layoutSizingHorizontal = "FILL";
    const userHeader = createAL("HORIZONTAL", 6, 0, 0);
    userHeader.counterAxisAlignItems = "CENTER";
    userHeader.appendChild(createTxt("👤 Iqbal Aditya", 12, "Semi Bold", rgb(147, 197, 253)));
    userHeader.appendChild(createTxt("03:28", 10, "Regular", T.textDim));
    userMsg.appendChild(userHeader);
    const userTxt = createTxt("Buatkan goroutine worker pool di Go dengan token bucket rate limiter 5,000 req/s dan AST unit test.", 13, "Regular", T.textWhite);
    userTxt.layoutSizingHorizontal = "FILL";
    userMsg.appendChild(userTxt);
    chatPane.appendChild(userMsg);

    // Message 2: AI Assistant Response Container
    const aiMsg = createAL("VERTICAL", 12, 16, 14, solid(22, 25, 32));
    aiMsg.cornerRadius = 12;
    aiMsg.strokes = solid(40, 46, 58);
    aiMsg.strokeWeight = 1;
    aiMsg.layoutSizingHorizontal = "FILL";

    const aiHeader = createAL("HORIZONTAL", 8, 0, 0);
    aiHeader.counterAxisAlignItems = "CENTER";
    const aiIcon = createAL("HORIZONTAL", 0, 4, 4, solid(99, 102, 241));
    aiIcon.cornerRadius = 4;
    aiIcon.appendChild(createTxt("⚡", 10, "Regular", T.textWhite));
    aiHeader.appendChild(aiIcon);
    aiHeader.appendChild(createTxt("Adith AI Engine", 12, "Bold", rgb(199, 210, 254)));
    aiHeader.appendChild(createTxt("Thinking Completed in 0.6s", 10, "Regular", T.accentEmerald));
    aiMsg.appendChild(aiHeader);

    // Collapsible Thinking Accordion Box
    const thinkBox = createAL("VERTICAL", 6, 12, 10, solid(13, 15, 18));
    thinkBox.cornerRadius = 8;
    thinkBox.strokes = solid(40, 46, 58);
    thinkBox.strokeWeight = 1;
    thinkBox.layoutSizingHorizontal = "FILL";
    thinkBox.appendChild(createTxt("🧠 Reasoning & AST Chain of Thought", 11, "Semi Bold", rgb(165, 180, 252)));
    thinkBox.appendChild(createTxt("• Tree-sitter AST: Membedah struct Worker & Pool pada context workspace", 11, "Regular", T.textMuted));
    thinkBox.appendChild(createTxt("• Qdrant Vector: 3 chunks kode terkait rate-limiting ditemukan di @limiter.go", 11, "Regular", T.textMuted));
    thinkBox.appendChild(createTxt("• Arsitektur: Menggunakan sync.WaitGroup + atomic backpressure queue", 11, "Regular", T.textMuted));
    aiMsg.appendChild(thinkBox);

    // AI Explanation Text
    const aiText = createTxt("Saya telah merancang implementasi WorkerPool dengan buffered channel dan rate limiter token bucket 5,000 req/s. Kode dan diff live siap Anda tinjau di Code Canvas sebelah kanan.", 13, "Regular", T.textWhite);
    aiText.layoutSizingHorizontal = "FILL";
    aiMsg.appendChild(aiText);

    // Quick Action Chips
    const chipsRow = createAL("HORIZONTAL", 6, 0, 0);
    const chip1 = createAL("HORIZONTAL", 4, 10, 6, solid(30, 35, 45));
    chip1.cornerRadius = 16;
    chip1.appendChild(createTxt("🔍 Jelaskan AST", 11, "Medium", T.textWhite));
    const chip2 = createAL("HORIZONTAL", 4, 10, 6, solid(30, 35, 45));
    chip2.cornerRadius = 16;
    chip2.appendChild(createTxt("🧪 Tambah Benchmark", 11, "Medium", T.textWhite));
    const chip3 = createAL("HORIZONTAL", 4, 10, 6, solid(30, 35, 45));
    chip3.cornerRadius = 16;
    chip3.appendChild(createTxt("⚡ Optimasi Memori", 11, "Medium", T.textWhite));
    chipsRow.appendChild(chip1);
    chipsRow.appendChild(chip2);
    chipsRow.appendChild(chip3);
    aiMsg.appendChild(chipsRow);

    chatPane.appendChild(aiMsg);

    // Chat Spacer
    const chatSpacer = figma.createFrame();
    chatSpacer.layoutGrow = 1;
    chatSpacer.fills = [];
    chatPane.appendChild(chatSpacer);

    // Prompt Input Bar at Bottom
    const inputContainer = createAL("VERTICAL", 8, 14, 12, solid(22, 25, 32));
    inputContainer.cornerRadius = 12;
    inputContainer.strokes = solid(55, 63, 80);
    inputContainer.strokeWeight = 1;
    inputContainer.layoutSizingHorizontal = "FILL";

    // Mentions tag row
    const tagRow = createAL("HORIZONTAL", 6, 0, 0);
    const tag1 = createAL("HORIZONTAL", 4, 8, 3, solid(30, 35, 45));
    tag1.cornerRadius = 4;
    tag1.appendChild(createTxt("@worker.go", 10, "Medium", rgb(147, 197, 253)));
    const tag2 = createAL("HORIZONTAL", 4, 8, 3, solid(30, 35, 45));
    tag2.cornerRadius = 4;
    tag2.appendChild(createTxt("@limiter.go", 10, "Medium", rgb(147, 197, 253)));
    tagRow.appendChild(tag1);
    tagRow.appendChild(tag2);
    inputContainer.appendChild(tagRow);

    const inputRow = createAL("HORIZONTAL", 8, 0, 0);
    inputRow.counterAxisAlignItems = "CENTER";
    inputRow.layoutSizingHorizontal = "FILL";
    const inputPlaceholder = createTxt("Tanya Adith AI atau ketik @ untuk konteks file...", 13, "Regular", T.textDim);
    inputPlaceholder.layoutGrow = 1;
    const btnSend = createAL("HORIZONTAL", 0, 10, 8, solid(99, 102, 241));
    btnSend.cornerRadius = 6;
    btnSend.appendChild(createTxt("Kirim ↵", 11, "Bold", rgb(255, 255, 255)));
    inputRow.appendChild(inputPlaceholder);
    inputRow.appendChild(btnSend);
    inputContainer.appendChild(inputRow);

    chatPane.appendChild(inputContainer);
    workspace.appendChild(chatPane);

    // ------------------------------------------
    // RIGHT PANE: CODE & CANVAS (900px)
    // ------------------------------------------
    const codePane = createAL("VERTICAL", 0, 0, 0, solid(13, 15, 18));
    codePane.resize(900, 856);
    codePane.layoutSizingHorizontal = "FIXED";
    codePane.layoutSizingVertical = "FILL";

    // File Tabs Bar
    const tabsBar = createAL("HORIZONTAL", 2, 16, 8, solid(17, 20, 26));
    tabsBar.resize(900, 42);
    tabsBar.layoutSizingHorizontal = "FIXED";
    tabsBar.strokes = solid(40, 46, 58);
    tabsBar.strokeWeight = 1;

    const tabActive = createAL("HORIZONTAL", 8, 12, 6, solid(22, 25, 32));
    tabActive.cornerRadius = 6;
    tabActive.strokes = solid(40, 46, 58);
    tabActive.strokeWeight = 1;
    tabActive.appendChild(createTxt("📄 worker_pool.go", 12, "Medium", T.textWhite));
    tabActive.appendChild(createTxt("×", 12, "Regular", T.textDim));

    const tab2 = createAL("HORIZONTAL", 8, 12, 6, []);
    tab2.appendChild(createTxt("📄 limiter.go", 12, "Regular", T.textMuted));

    const tab3 = createAL("HORIZONTAL", 8, 12, 6, []);
    tab3.appendChild(createTxt("🧪 pool_test.go", 12, "Regular", T.textMuted));

    const tab4 = createAL("HORIZONTAL", 8, 12, 6, []);
    tab4.appendChild(createTxt("👁️ Live Preview", 12, "Regular", T.accentEmerald));

    tabsBar.appendChild(tabActive);
    tabsBar.appendChild(tab2);
    tabsBar.appendChild(tab3);
    tabsBar.appendChild(tab4);
    codePane.appendChild(tabsBar);

    // Diff Action Bar (Green Banner)
    const diffBar = createAL("HORIZONTAL", 12, 20, 10, solid(16, 185, 129, 0.15));
    diffBar.resize(900, 48);
    diffBar.layoutSizingHorizontal = "FIXED";
    diffBar.counterAxisAlignItems = "CENTER";
    diffBar.strokes = solid(16, 185, 129, 0.4);
    diffBar.strokeWeight = 1;

    const diffBadge = createAL("HORIZONTAL", 0, 8, 4, solid(16, 185, 129));
    diffBadge.cornerRadius = 4;
    diffBadge.appendChild(createTxt("+28 / -4 lines", 11, "Bold", rgb(255, 255, 255)));
    diffBar.appendChild(diffBadge);

    diffBar.appendChild(createTxt("Perubahan otomatis siap di-merge ke codebase", 12, "Medium", rgb(209, 250, 229)));

    const diffSpacer = figma.createFrame();
    diffSpacer.layoutGrow = 1;
    diffSpacer.fills = [];
    diffBar.appendChild(diffSpacer);

    const btnAccept = createAL("HORIZONTAL", 4, 12, 6, solid(16, 185, 129));
    btnAccept.cornerRadius = 6;
    btnAccept.appendChild(createTxt("✓ Accept All", 11, "Bold", rgb(255, 255, 255)));

    const btnDiscard = createAL("HORIZONTAL", 4, 12, 6, solid(30, 35, 45));
    btnDiscard.cornerRadius = 6;
    btnDiscard.appendChild(createTxt("✗ Discard", 11, "Medium", T.textMuted));

    diffBar.appendChild(btnAccept);
    diffBar.appendChild(btnDiscard);
    codePane.appendChild(diffBar);

    // Monaco Code Editor Body (Simulated)
    const codeEditor = createAL("VERTICAL", 4, 24, 16, solid(13, 15, 18));
    codeEditor.resize(900, 520);
    codeEditor.layoutSizingHorizontal = "FIXED";
    codeEditor.layoutGrow = 1;

    const codeLines = [
        { num: "01", text: "package worker", color: rgb(244, 63, 94) },
        { num: "02", text: "", color: T.textWhite },
        { num: "03", text: "import (", color: rgb(244, 63, 94) },
        { num: "04", text: '    "context"', color: rgb(251, 191, 36) },
        { num: "05", text: '    "sync"', color: rgb(251, 191, 36) },
        { num: "06", text: '    "golang.org/x/time/rate"', color: rgb(251, 191, 36) },
        { num: "07", text: ")", color: rgb(244, 63, 94) },
        { num: "08", text: "", color: T.textWhite },
        { num: "09", text: "// WorkerPool coordinates high-throughput execution with rate limiting", color: rgb(100, 116, 139) },
        { num: "10", text: "type WorkerPool struct {", color: rgb(99, 102, 241) },
        { num: "11", text: "    tasks     chan Task", color: T.textWhite },
        { num: "12", text: "    limiter   *rate.Limiter", color: T.textWhite },
        { num: "13", text: "    wg        sync.WaitGroup", color: T.textWhite },
        { num: "14", text: "}", color: rgb(99, 102, 241) },
        { num: "15", text: "", color: T.textWhite },
        { num: "16", text: "+ func NewWorkerPool(workers int, rps rate.Limit) *WorkerPool {", color: T.accentEmerald, isDiff: true },
        { num: "17", text: "+     return &WorkerPool{", color: T.accentEmerald, isDiff: true },
        { num: "18", text: "+         tasks:   make(chan Task, workers*8),", color: T.accentEmerald, isDiff: true },
        { num: "19", text: "+         limiter: rate.NewLimiter(rps, int(rps)),", color: T.accentEmerald, isDiff: true },
        { num: "20", text: "+     }", color: T.accentEmerald, isDiff: true },
        { num: "21", text: "+ }", color: T.accentEmerald, isDiff: true },
    ];

    for (const l of codeLines) {
        const lineRow = createAL("HORIZONTAL", 16, 8, 2, l.isDiff ? solid(16, 185, 129, 0.12) : []);
        lineRow.cornerRadius = 4;
        lineRow.layoutSizingHorizontal = "FILL";
        lineRow.appendChild(createTxt(l.num, 12, "Regular", T.textDim));
        lineRow.appendChild(createTxt(l.text, 13, l.isDiff ? "Semi Bold" : "Regular", l.color));
        codeEditor.appendChild(lineRow);
    }
    codePane.appendChild(codeEditor);

    // Bottom Terminal & AST Drawer (220px)
    const terminalDrawer = createAL("VERTICAL", 8, 20, 12, solid(17, 20, 26));
    terminalDrawer.resize(900, 210);
    terminalDrawer.layoutSizingHorizontal = "FIXED";
    terminalDrawer.strokes = solid(40, 46, 58);
    terminalDrawer.strokeWeight = 1;

    const termHead = createAL("HORIZONTAL", 12, 0, 0);
    termHead.appendChild(createTxt("⚙️ TERMINAL SANDBOX", 11, "Bold", T.textWhite));
    termHead.appendChild(createTxt("📊 AST DIAGNOSTICS", 11, "Medium", T.textDim));
    termHead.appendChild(createTxt("⚡ QDRANT VECTORS", 11, "Medium", T.textDim));
    terminalDrawer.appendChild(termHead);

    const termBox = createAL("VERTICAL", 4, 12, 10, solid(13, 15, 18));
    termBox.cornerRadius = 6;
    termBox.layoutSizingHorizontal = "FILL";
    termBox.appendChild(createTxt("$ go test -v -bench=. ./internal/worker/...", 11, "Medium", T.textMuted));
    termBox.appendChild(createTxt("=== RUN   TestWorkerPool_Throughput", 11, "Regular", T.textDim));
    termBox.appendChild(createTxt("--- PASS: TestWorkerPool_Throughput (0.04s) - 5,280 ops/sec sustained", 11, "Semi Bold", T.accentEmerald));
    termBox.appendChild(createTxt("BenchmarkWorkerPool-16    1000000    1840 ns/op    120 B/op", 11, "Regular", T.accentAmber));
    termBox.appendChild(createTxt("PASS: ok  iqbal-portfolio/internal/worker  0.182s", 11, "Bold", T.accentEmerald));
    terminalDrawer.appendChild(termBox);

    codePane.appendChild(terminalDrawer);
    workspace.appendChild(codePane);
    root.appendChild(workspace);

    // Zoom view into newly created frame
    figma.viewport.scrollAndZoomIntoView([root]);
    console.log("✅ ADITH AI STUDIO frame generated successfully at (" + startX + ", 0)!");
})();
