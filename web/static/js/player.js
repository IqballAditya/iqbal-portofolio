(function () {
    const rawData = document.getElementById('player-milestones');
    let milestones = [];
    if (rawData) {
        try {
            milestones = JSON.parse(rawData.textContent || '[]');
        } catch (e) {
            console.error('Failed to parse milestones JSON:', e);
        }
    }

    let isPlaying = true;
    let currentEpIndex = 0;
    let speed = 1.0;
    let progress = 0.0;
    let logInterval = null;
    let stepCount = 0;

    const playPauseBtn = document.getElementById('play-pause-btn');
    const iconPlay = document.getElementById('icon-play');
    const iconPause = document.getElementById('icon-pause');
    const rewindBtn = document.getElementById('rewind-btn');
    const speedButtons = document.querySelectorAll('.speed-btn');
    const scrubberTrack = document.getElementById('scrubber-track');
    const scrubberProgress = document.getElementById('scrubber-progress');
    const scrubberThumb = document.getElementById('scrubber-thumb');
    const terminalLines = document.getElementById('terminal-lines');
    const terminalOutput = document.getElementById('terminal-output');
    const currentEpTitleEl = document.getElementById('player-current-episode-title');
    const snippetBodyEl = document.getElementById('code-snippet-body');
    const timeDisplayEl = document.getElementById('player-time-display');
    const drawerToggleBtn = document.getElementById('drawer-toggle-btn');
    const episodesDrawerBtn = document.getElementById('episodes-drawer-btn');
    const drawerCloseBtn = document.getElementById('drawer-close-btn');
    const episodesDrawer = document.getElementById('episodes-drawer');
    const chapterTabs = document.querySelectorAll('.chapter-tab');
    const drawerItems = document.querySelectorAll('.drawer-item');

    function getLogTemplates(epTitle) {
        return [
            { text: "[INIT] go runtime 1.22.4 (darwin/arm64) GOMAXPROCS=16", cls: "log-info" },
            { text: "[CONFIG] worker_pool=128 channel_buffer=4096 sync_strategy=waitgroup", cls: "log-cyan" },
            { text: "[BOOT] spawning goroutine cluster for " + epTitle, cls: "log-green" },
            { text: "[DAG] resolving dependency node graph: 14 topological stages identified", cls: "log-purple" },
            { text: "[MEM] heap_alloc=18.4MB sys=32.1MB next_gc=36.8MB gc_pause=140µs", cls: "log-dim" },
            { text: "[BENCHMARK] BenchmarkDAGWalkConcurrent-16  500000  2410 ns/op  320 B/op", cls: "log-yellow" },
            { text: "[RAFT] leader term=4 committed_index=192408 quorum_active=true", cls: "log-green" },
            { text: "[WORKER-07] task_batch=512 dispatched, latency=840µs status=SUCCESS", cls: "log-cyan" },
            { text: "[CHAOS] injected network jitter: latency spiked to 3.2ms, recovered in 12ms", cls: "log-yellow" },
            { text: "[eBPF] kprobe/tcp_v4_connect captured 4,120 socket events with zero dropped frames", cls: "log-purple" },
            { text: "[PASS] TestByzantineFaultTolerance (0.04s) - 0 errors, 100% consensus reached", cls: "log-green" },
            { text: "[HEALTH] cluster status: HEALTHY (p50=0.8ms, p99=1.2ms, throughput=250k ops/sec)", cls: "log-green" }
        ];
    }

    function appendLine(text, cls) {
        if (!terminalLines) return;
        const line = document.createElement('div');
        line.className = 'terminal-line ' + (cls || '');
        const time = new Date().toISOString().substring(11, 19) + '.' + Math.floor(Math.random() * 900 + 100);
        line.innerHTML = '<span class="log-timestamp">' + time + '</span> ' + escapeHtml(text);
        terminalLines.appendChild(line);

        while (terminalLines.children.length > 80) {
            terminalLines.removeChild(terminalLines.firstChild);
        }

        if (terminalOutput) {
            terminalOutput.scrollTop = terminalOutput.scrollHeight;
        }
    }

    function escapeHtml(s) {
        return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    }

    function updateUIForEpisode(index) {
        if (!milestones || milestones.length === 0) return;
        if (index < 0 || index >= milestones.length) index = 0;
        currentEpIndex = index;
        const ep = milestones[currentEpIndex];

        if (currentEpTitleEl && ep) {
            currentEpTitleEl.textContent = ep.title;
        }
        if (snippetBodyEl && ep) {
            snippetBodyEl.textContent = ep.code_snippet || '// No code snippet for this episode';
        }

        chapterTabs.forEach((tab, i) => {
            if (i === currentEpIndex) {
                tab.classList.add('active');
                tab.setAttribute('aria-selected', 'true');
            } else {
                tab.classList.remove('active');
                tab.setAttribute('aria-selected', 'false');
            }
        });

        drawerItems.forEach((item, i) => {
            if (i === currentEpIndex) {
                item.classList.add('active');
            } else {
                item.classList.remove('active');
            }
        });

        appendLine("--- SWITCHED TO CHAPTER: " + (ep ? ep.title : "Episode") + " ---", "log-header");
    }

    function setSpeed(newSpeed) {
        speed = newSpeed;
        speedButtons.forEach(btn => {
            if (parseFloat(btn.dataset.speed) === newSpeed) {
                btn.classList.add('active');
            } else {
                btn.classList.remove('active');
            }
        });
        if (isPlaying) {
            stopLoop();
            startLoop();
        }
    }

    function tick() {
        stepCount++;
        progress += (0.4 * speed);
        if (progress > 100) {
            progress = 0;
            if (milestones.length > 1) {
                updateUIForEpisode((currentEpIndex + 1) % milestones.length);
            }
        }

        if (scrubberProgress) scrubberProgress.style.width = progress + '%';
        if (scrubberThumb) scrubberThumb.style.left = progress + '%';
        if (scrubberTrack) scrubberTrack.setAttribute('aria-valuenow', Math.round(progress));

        const totalSec = 2700;
        const curSec = Math.floor((progress / 100) * totalSec);
        const curM = String(Math.floor(curSec / 60)).padStart(2, '0');
        const curS = String(curSec % 60).padStart(2, '0');
        if (timeDisplayEl) {
            timeDisplayEl.textContent = curM + ':' + curS + ' / 45:00';
        }

        const ep = milestones[currentEpIndex];
        const templates = getLogTemplates(ep ? ep.title : "Adith AI Engine");
        const entry = templates[stepCount % templates.length];
        appendLine(entry.text, entry.cls);
    }

    function startLoop() {
        if (logInterval) clearInterval(logInterval);
        const intervalMs = Math.max(180, Math.floor(600 / speed));
        logInterval = setInterval(tick, intervalMs);
        if (iconPlay) iconPlay.style.display = 'none';
        if (iconPause) iconPause.style.display = 'block';
    }

    function stopLoop() {
        if (logInterval) {
            clearInterval(logInterval);
            logInterval = null;
        }
        if (iconPlay) iconPlay.style.display = 'block';
        if (iconPause) iconPause.style.display = 'none';
    }

    function togglePlay() {
        isPlaying = !isPlaying;
        if (isPlaying) {
            startLoop();
        } else {
            stopLoop();
        }
    }

    if (playPauseBtn) {
        playPauseBtn.addEventListener('click', togglePlay);
    }

    if (rewindBtn) {
        rewindBtn.addEventListener('click', function () {
            progress = 0;
            if (terminalLines) terminalLines.innerHTML = '';
            appendLine("=== DEMO RESTARTED: " + (milestones[currentEpIndex] ? milestones[currentEpIndex].title : "") + " ===", "log-header");
            if (scrubberProgress) scrubberProgress.style.width = '0%';
            if (scrubberThumb) scrubberThumb.style.left = '0%';
        });
    }

    speedButtons.forEach(btn => {
        btn.addEventListener('click', function () {
            setSpeed(parseFloat(this.dataset.speed));
        });
    });

    chapterTabs.forEach(tab => {
        tab.addEventListener('click', function () {
            const idx = parseInt(this.dataset.episodeIndex, 10);
            updateUIForEpisode(idx);
        });
    });

    drawerItems.forEach(item => {
        item.addEventListener('click', function () {
            const idx = parseInt(this.dataset.episodeIndex, 10);
            updateUIForEpisode(idx);
            if (episodesDrawer) {
                episodesDrawer.classList.remove('open');
                episodesDrawer.setAttribute('aria-hidden', 'true');
            }
        });
    });

    if (scrubberTrack) {
        scrubberTrack.addEventListener('click', function (e) {
            const rect = scrubberTrack.getBoundingClientRect();
            const clickX = e.clientX - rect.left;
            progress = Math.max(0, Math.min(100, (clickX / rect.width) * 100));
            if (scrubberProgress) scrubberProgress.style.width = progress + '%';
            if (scrubberThumb) scrubberThumb.style.left = progress + '%';
        });
    }

    function toggleDrawer() {
        if (!episodesDrawer) return;
        const isOpen = episodesDrawer.classList.contains('open');
        if (isOpen) {
            episodesDrawer.classList.remove('open');
            episodesDrawer.setAttribute('aria-hidden', 'true');
        } else {
            episodesDrawer.classList.add('open');
            episodesDrawer.setAttribute('aria-hidden', 'false');
        }
    }

    if (drawerToggleBtn) drawerToggleBtn.addEventListener('click', toggleDrawer);
    if (episodesDrawerBtn) episodesDrawerBtn.addEventListener('click', toggleDrawer);
    if (drawerCloseBtn) drawerCloseBtn.addEventListener('click', toggleDrawer);

    window.addEventListener('keydown', function (e) {
        if (e.target && (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA')) return;
        if (e.code === 'Space') {
            e.preventDefault();
            togglePlay();
        }
    });

    if (milestones.length > 0) {
        updateUIForEpisode(0);
    }
    startLoop();
})();
