const querySelectorAll = (selector, element = document) =>
  Array.prototype.slice.call(element.querySelectorAll(selector));

querySelectorAll('.js-guest').forEach((guest) => {
  const mealOptions = guest.querySelector('.js-meal-options');

  querySelectorAll('.js-reception-toggle', guest).forEach((toggle) => {
    toggle.addEventListener('change', () => {
      if (toggle.checked) {
        mealOptions.classList.remove('hidden');
      } else {
        mealOptions.classList.add('hidden');
      }
    });
  });
});
