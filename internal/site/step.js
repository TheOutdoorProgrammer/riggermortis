// SPDX-License-Identifier: Apache-2.0
// Progressive enhancement only. The radios already work without this file;
// all it adds is arrows, keyboard control and autoplay.
(function () {
  var radios = Array.prototype.slice.call(
    document.querySelectorAll('.stepper input[type=radio]'));
  if (radios.length < 2) return;

  function current() {
    for (var i = 0; i < radios.length; i++) if (radios[i].checked) return i;
    return 0;
  }
  function go(i) { radios[(i + radios.length) % radios.length].checked = true; }

  var prev = document.getElementById('prev-btn');
  var next = document.getElementById('next-btn');
  if (prev) prev.addEventListener('click', function () { go(current() - 1); });
  if (next) next.addEventListener('click', function () { go(current() + 1); });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'ArrowLeft') { go(current() - 1); }
    else if (e.key === 'ArrowRight') { go(current() + 1); }
  });
})();
