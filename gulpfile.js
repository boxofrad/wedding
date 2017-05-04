const gulp = require('gulp');
const sass = require('gulp-sass');
const rev = require('gulp-rev');
const revReplace = require('gulp-rev-replace');
const filter = require('gulp-filter');

gulp.task('default', () => {
  const imageFilter = filter('**/*.jpg', { restore: true });
  const sassFilter = filter('**/*.scss', { restore: true });
  const cssFilter = filter('**/*.css', { restore: true });

  gulp.src(['./assets/images/*', './assets/stylesheets/application.scss'])
      // fingerprint images
      .pipe(imageFilter)
      .pipe(rev())
      .pipe(imageFilter.restore)

      // compile SASS
      .pipe(sassFilter)
      .pipe(sass({ outputStyle: 'compressed' }).on('error', sass.logError))
      .pipe(sassFilter.restore)

      // replace image references with fingerprinted file names
      .pipe(revReplace())

      // fingerprint the CSS file now that the image references have been rewritten
      .pipe(cssFilter)
      .pipe(rev())
      .pipe(cssFilter.restore)
      .pipe(gulp.dest('./public'))

      // save out a manifest
      .pipe(rev.manifest())
      .pipe(gulp.dest('./public'));
});

gulp.task('watch', () => {
  gulp.watch('./assets/**/*', ['default']);
});
