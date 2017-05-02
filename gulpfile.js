const gulp = require('gulp');
const sass = require('gulp-sass');
const rev = require('gulp-rev');

gulp.task('default', () => {
  gulp.src('./assets/images/*')
      .pipe(rev())
      .pipe(gulp.dest('./public'));

  gulp.src('./assets/stylesheets/application.scss')
      .pipe(sass().on('error', sass.logError))
      .pipe(rev())
      .pipe(gulp.dest('./public'));
});
