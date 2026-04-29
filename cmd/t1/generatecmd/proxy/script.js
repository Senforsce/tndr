(function() {
  let tndr_reloadSrc = window.t1_reloadSrc || new EventSource("/_t1/reload/events");
  tndr_reloadSrc.onmessage = (event) => {
    if (event && event.data === "reload") {
      window.location.reload();
    }
  };
  window.t1_reloadSrc = tndr_reloadSrc;
  window.onbeforeunload = () => window.t1_reloadSrc.close();
})();
