#ifndef GTK_C_CALLBACKS
#define GTK_C_CALLBACKS

#include <gtk/gtk.h>

extern void gtkgo_quit_c();
extern gboolean gtkgo_gl_area_render_c(GtkGLArea *area, GdkGLContext *context);
extern void gtkgo_gl_area_realize_c(GtkGLArea *area);

extern gboolean queue_render_idle(gpointer user_data); 

#endif