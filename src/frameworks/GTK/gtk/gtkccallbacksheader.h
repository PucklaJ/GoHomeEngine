#ifndef GTK_C_CALLBACKS
#define GTK_C_CALLBACKS

#include <gtk/gtk.h>

extern void gtkgo_quit_c();
extern gboolean gtkgo_gl_area_render_c(GtkGLArea *area, GdkGLContext *context);
extern void gtkgo_gl_area_realize_c(GtkGLArea *area);
extern gboolean gtkgo_gl_area_key_press_c(GtkWidget* widget,GdkEvent* event,gpointer user_data);
extern gboolean gtkgo_gl_area_key_release_c(GtkWidget* widget, GdkEvent* event,gpointer user_data);
extern gboolean gtkgo_gl_area_button_press_c(GtkWidget *widget, GdkEvent *event, gpointer user_data);
extern gboolean gtkgo_gl_area_button_release_c(GtkWidget *widget, GdkEvent *event, gpointer user_data);
extern gboolean gtkgo_gl_area_motion_notify_c(GtkWidget *widget, GdkEvent *event, gpointer user_data);
extern gboolean gtkgo_gl_area_scroll_c(GtkWidget *widget, GdkEvent *event, gpointer user_data);
extern void gtkgo_button_signal_c(GtkButton *button, gpointer user_data);
extern void gtkgo_widget_size_allocate_signal_c(GtkWidget* widget, GdkRectangle* allocation, gpointer user_data);
extern void gtkgo_menu_item_signal_c(GtkMenuItem* menuItem, gpointer user_data);
extern void gtkgo_widget_event_signal_c(GtkWidget* widget, GdkEvent* event, gpointer user_data);
extern void gtkgo_list_box_row_selected_signal_c(GtkListBox* listBox, GtkListBoxRow* listBoxRow, gpointer user_data);
extern void gtkgo_tool_button_signal_c(GtkToolButton* toolButton, gpointer user_data);

extern gboolean queue_render_idle(gpointer user_data); 

#endif