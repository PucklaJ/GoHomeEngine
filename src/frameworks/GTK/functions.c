#include "includes.h"
#include "gtkccallbacksheader.h"

GtkWindow* Window = NULL;
GtkGLArea* GLarea = NULL;
char* ErrorString = NULL;

void initialise(int args,char** argv)
{
	ErrorString = (char*)malloc(150);
	gtk_init(&args,&argv);
}

int createWindow(unsigned int width, unsigned int height, const char* title)
{
	Window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
	gtk_widget_set_size_request(GTK_WIDGET(Window),width,height);
	gtk_window_set_title(Window,title);
	gtk_widget_set_events(GTK_WIDGET(Window), GDK_POINTER_MOTION_MASK|GDK_SCROLL_MASK);

	g_signal_connect(GTK_WIDGET(Window),"delete-event",G_CALLBACK(gtkgo_quit_c),NULL);

	GLarea = gtk_gl_area_new();
	g_signal_connect(GTK_WIDGET(GLarea),"render",G_CALLBACK(gtkgo_gl_area_render_c),NULL);
	g_signal_connect(GTK_WIDGET(GLarea),"realize",G_CALLBACK(gtkgo_gl_area_realize_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"key-press-event",G_CALLBACK(gtkgo_gl_area_key_press_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"key-release-event",G_CALLBACK(gtkgo_gl_area_key_release_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"button-press-event",G_CALLBACK(gtkgo_gl_area_button_press_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"button-release-event",G_CALLBACK(gtkgo_gl_area_button_release_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"motion-notify-event",G_CALLBACK(gtkgo_gl_area_motion_notify_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"scroll-event",G_CALLBACK(gtkgo_gl_area_scroll_c),NULL);
	gtk_gl_area_set_has_depth_buffer(GLarea,TRUE);
	gtk_container_add(GTK_CONTAINER(Window),GTK_WIDGET(GLarea));

	gdk_threads_add_idle(queue_render_idle,GLarea);

	gtk_widget_show_all(GTK_WIDGET(Window));
	return 1;
}

void windowGetSize(float* width, float* height)
{
	if(width == NULL || height == NULL)
	{
		return;
	}

	GtkAllocation* alloc = g_new(GtkAllocation, 1);
    gtk_widget_get_allocation(GTK_WIDGET(Window), alloc);
    if(alloc != NULL)
    {
    	*width = alloc->width;
    	*height = alloc->height;    	
    	g_free(alloc);
    }
    else
    {
    	*width = -1.0;
    	*height = -1.0;
    }
}

void loop()
{
	gtk_main();
}