const querySelectorAll = (selector, element = document) =>
  Array.prototype.slice.call(element.querySelectorAll(selector));

querySelectorAll('.js-guest').forEach((guest) => {
  const attendingFields = guest.querySelector('.js-attending-fields');

  querySelectorAll('.js-attending-toggle', guest).forEach((toggle) => {
    toggle.addEventListener('click', () => {
      if (toggle.value === '1') {
        attendingFields.classList.remove('hidden');
      } else {
        attendingFields.classList.add('hidden');
      }
    });
  });
});
