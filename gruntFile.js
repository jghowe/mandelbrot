module.exports = function (grunt) {

    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-html2js');

    // Default task
    grunt.registerTask('default', ['build']);
    grunt.registerTask('install', ['clean', 'copy', 'html2js', 'concat']);

    // Print a timestamp (usefule for when watching)
    grunt.registerTask('timestamp', function () {
        grunt.log.subhead(Date());
    });

    // Project configuration.
    grunt.initConfig({
        distdir: 'static',
        pkg: grunt.file.readJSON('package.json'),
        src: {
            css: ['client/css/**/*.css'],
            js: ['client/js/**/*.js'],
            jsTpl: ['<%= distdir %>/templates/**/*.js'],
            html: ['client/*.html'],
            tpl: {
                client: ['client/js/**/*.template.html']
            }
        },
        clean: ['<%= distdir %>/*'],
        concat: {
            dist: {
                options: {
                    sourceMap: true
                },
                src: ['<%= src.js %>', '<%= src.jsTpl %>'],
                dest: '<%= distdir %>/js/<%= pkg.name %>.js'
            },
            jquery: {
                options: {
                    sourceMap: true
                },
                src: ['bower_components/jquery/dist/jquery.min.js'],
                dest: '<%= distdir %>/js/jquery.js'
            },
            angular: {
                options: {
                    sourceMap: true
                },
                src: [
                    'bower_components/angular/angular.min.js',
                    'bower_components/angular-resource/angular-resource.min.js'
                ],
                dest: '<%= distdir %>/js/angular.js'
            },
            angularui: {
                options: {
                    sourceMap: true
                },
                src: [
                    'bower_components/angular-ui-router/release/angular-ui-router.min.js'
                ],
                dest: '<%= distdir %>/js/angular-ui.js'
            },
            components: {
                options: {
                    sourceMap: true
                },
                src: [
                    'bower_components/spin.js/spin.js',
                    'bower_components/angular-spinner/angular-spinner.js',
                    'bower_components/hamsterjs/hamster.js',
                    'bower_components/angular-mousewheel/mousewheel.js',
                    'bower_components/angular-pan-zoom/release/panzoom.min.js'
                ],
                dest: '<%= distdir %>/js/components.js'
            }
        },
        copy: {
            images: {
                files: [{
                    dest: '<%= distdir %>/img/',
                    src: '*',
                    expand: true,
                    cwd: 'client/img'
                }]
            },
            css: {
                files: [{
                    dest: '<%= distdir %>/css/',
                    src: '*.css',
                    expand: true,
                    cwd: 'bower_components/angular-pan-zoom/release'
                }, {
                    dest: '<%= distdir %>/css/',
                    src: ['foundation.css', 'foundation.css.map'],
                    expand: true,
                    cwd: 'bower_components/foundation/css'
                }, {
                    dest: '<%= distdir %>/css/',
                    src: '*.css',
                    expand: true,
                    cwd: 'client/css'
                }]
            },
            index: {
                files: [{
                    dest: '<%= distdir %>/index.html',
                    src: 'client/index.html'
                }]
            }
        },
        html2js: {
            client: {
                options: {
                    base: 'client/'
                },
                src: ['<%= src.tpl.client %>'],
                dest: '<%= distdir %>/templates/app.js',
                module: 'templates.app'
            }
        },
        watch: {
            build: {
                files: ['<%= src.css %>', '<%= src.js %>', '<%= src.html %>', '<%= src.tpl.client %>'],
                tasks: ['copy:css', 'copy:index', 'html2js:client', 'concat:dist']
            }
        }
    });
}
