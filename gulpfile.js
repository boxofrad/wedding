const gulp = require('gulp');
const babel = require('gulp-babel');
const uglify = require('gulp-uglify');
const sourcemaps = require('gulp-sourcemaps');
const flatten = require('gulp-flatten');
const sass = require('gulp-sass');
const rev = require('gulp-rev');
const revReplace = require('gulp-rev-replace');
const filter = require('gulp-filter');

gulp.task('default', () => {
  const imageFilter = filter('**/*.jpg', { restore: true });
  const sassFilter = filter('**/*.scss', { restore: true });
  const cssFilter = filter('**/*.css', { restore: true });
  const jsFilter = filter('**/*.js', { restore: true });

  gulp.src(['./assets/{images,javascripts}/*', './assets/stylesheets/application.scss'])
      // transpile JS
      .pipe(jsFilter)
      .pipe(sourcemaps.init())
      .pipe(babel({ presets: ['es2015'] }))
      .pipe(uglify())
      .pipe(rev())
      .pipe(flatten())
      .pipe(sourcemaps.write('.'))
      .pipe(gulp.dest('./static'))
      .pipe(jsFilter.restore)

      // fingerprint images
      .pipe(imageFilter)
      .pipe(rev())
      .pipe(imageFilter.restore)

      // compile SASS
      .pipe(sassFilter)
      .pipe(sass({ outputStyle: 'compressed' }).on('error', sass.logError))
      .pipe(sassFilter.restore)

      // replace image references with fingerprinted file names
      .pipe(revReplace({ prefix: '/static/' }))

      // fingerprint the CSS file now that the image references have been rewritten
      .pipe(cssFilter)
      .pipe(rev())
      .pipe(cssFilter.restore)
      .pipe(gulp.dest('./static'))

      // save out a manifest
      .pipe(rev.manifest())
      .pipe(gulp.dest('./static'));
});

gulp.task('watch', () => {
  gulp.watch('./assets/**/*', ['default']);
});
