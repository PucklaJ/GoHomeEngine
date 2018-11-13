#include "gtkccallbacksheader.h"
#include "_cgo_export.h"
#include "includes.h"

void gtkgo_quit_c()
{
	gtkgo_quit();
}

gboolean gtkgo_gl_area_render_c(GtkGLArea *area, GdkGLContext *context)
{
	gtkgo_gl_area_render(area,context);
	return TRUE;
}

void gtkgo_gl_area_realize_c(GtkGLArea *area)
{
	int err = 0;

	gtk_gl_area_make_current(area);

	if(gtk_gl_area_get_error(area) != NULL)
	{
		ErrorString = "Couldn't make context current";
		err = 1;
	}
	else
	{
		gtk_gl_area_set_auto_render(area,FALSE);
	}

	gtkgo_gl_area_realize(area,err);

	if(err != 0)
	{
		// gtk_gl_area_queue_render(area);
	}
}

gboolean queue_render_idle(gpointer user_data)
{
	GtkGLArea* area = user_data;
	gtk_widget_queue_draw(GTK_WIDGET(area));
	gtk_gl_area_queue_render(area);

	return TRUE;
}

gboolean gtkgo_gl_area_key_press_c(GtkWidget* widget,GdkEvent* event,gpointer user_data)
{
	gtkgo_gl_area_key_press(widget,(GdkEventKey*)event);
	return TRUE;
}

gboolean gtkgo_gl_area_key_release_c(GtkWidget* widget, GdkEvent* event,gpointer user_data)
{
	gtkgo_gl_area_key_release(widget,(GdkEventKey*)event);
	return TRUE;
}

gboolean gtkgo_gl_area_button_press_c(GtkWidget *widget, GdkEvent *event, gpointer user_data)
{
	gtkgo_gl_area_button_press(widget,(GdkEventButton*)event);
	return TRUE;
}

gboolean gtkgo_gl_area_button_release_c(GtkWidget *widget, GdkEvent *event, gpointer user_data)
{
	gtkgo_gl_area_button_release(widget,(GdkEventButton*)event);
	return TRUE;
}

gboolean gtkgo_gl_area_motion_notify_c(GtkWidget *widget, GdkEvent *event, gpointer user_data)
{
	gtkgo_gl_area_motion_notify(widget,(GdkEventMotion*)event);
	return TRUE;
}

gboolean gtkgo_gl_area_scroll_c(GtkWidget *widget, GdkEvent *event, gpointer user_data)
{
	gtkgo_gl_area_scroll(widget,(GdkEventScroll*)event,event);
	return TRUE;
}

void gtkgo_button_signal_c(GtkButton *button, gpointer user_data)
{
    ButtonSignalUserData* bsud = (ButtonSignalUserData*)user_data;
    gtkgo_button_signal(button,bsud->id,bsud->signal);
}

void gtkgo_widget_size_allocate_signal_c(GtkWidget* widget, GdkRectangle* allocation, gpointer user_data)
{
	WidgetSignalUserData* wsud = (WidgetSignalUserData*)user_data;
	gtkgo_widget_signal(widget,wsud->name,wsud->signal);
}

