let b = document.getElementById("bar"),
    st = document.getElementById("st"),
    dn = document.getElementById("dn"),
    l = document.getElementById("lnk"),
    dr = document.getElementById("drop"),
    bk = document.getElementById("back"),
    la = 0,
    tm = 0,
    done = false;

function up(f) {
    if (!f || done) return;
    let x = new XMLHttpRequest();
    x.open("POST", "/upload");
    x.upload.onprogress = (e) => {
        if (e.lengthComputable) {
            b.style.width = (e.loaded / e.total) * 100 + "%";
            st.style.display = "block";
            let n = Date.now(),
                sp = la ? (e.loaded - la) / ((n - tm) / 1000) : 0;
            la = e.loaded;
            tm = n;
            let F = (b) =>
                b >= 1e9
                    ? (b / 1e9).toFixed(1) + "GB"
                    : b >= 1e6
                      ? (b / 1e6).toFixed(1) + "MB"
                      : b >= 1e3
                        ? (b / 1e3).toFixed(1) + "KB"
                        : b + "B";
            st.textContent = F(e.loaded) + " / " + F(e.total) + " - " + F(sp) + "/s";
        }
    };
    x.onload = () => {
        done = true;
        dr.style.display = "none";
        b.style.width = "0";
        st.style.display = "none";
        l.href = l.textContent = window.location + x.responseText;
        dn.style.display = "block";
        bk.style.display = "block";
    };
    dr.style.display = "none";
    x.send(f);
}

bk.addEventListener("click", (e) => {
    e.stopPropagation();
    done = false;
    dr.style.display = "flex";
    dn.style.display = "none";
    bk.style.display = "none";
    b.style.width = "0";
    st.style.display = "none";
    la = 0;
    tm = 0;
});
document.addEventListener("dragover", (e) => e.preventDefault());
document.addEventListener("drop", (e) => {
    e.preventDefault();
    if (done) return;
    up(e.dataTransfer.files[0]);
});
document.addEventListener("click", () => {
    if (done) return;
    let i = document.createElement("input");
    i.type = "file";
    i.onchange = () => up(i.files[0]);
    i.click();
});