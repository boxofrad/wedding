const querySelectorAll = (selector, element = document) =>
  Array.prototype.slice.call(element.querySelectorAll(selector));

querySelectorAll('.js-guest').forEach((guest) => {
  const receptionFields = guest.querySelector('.js-reception-fields');

  querySelectorAll('.js-reception-toggle', guest).forEach((toggle) => {
    toggle.addEventListener('click', () => {
      if (toggle.value === '1') {
        receptionFields.classList.remove('hidden');
      } else {
        receptionFields.classList.add('hidden');
      }
    });
  });
});
