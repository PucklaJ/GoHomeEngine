#include "includes.h"
#include "gtkccallbacksheader.h"

GtkWindow* window = NULL;
GtkGLArea* glarea = NULL;
char* ErrorString = NULL;

void initialise(int args,char** argv)
{
	ErrorString = (char*)malloc(150);
	gtk_init(&args,&argv);
}

int createWindow(unsigned int width, unsigned int height, const char* title)
{
	window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
	gtk_widget_set_size_request(GTK_WIDGET(window),width,height);
	gtk_window_set_title(window,title);

	g_signal_connect(GTK_WIDGET(window),"delete-event",G_CALLBACK(gtkgo_quit_c),NULL);

	glarea = gtk_gl_area_new();
	g_signal_connect(GTK_WIDGET(glarea),"render",G_CALLBACK(gtkgo_gl_area_render_c),NULL);
	g_signal_connect(GTK_WIDGET(glarea),"realize",G_CALLBACK(gtkgo_gl_area_realize_c),NULL);
	gtk_gl_area_set_has_depth_buffer(glarea,TRUE);
	gtk_container_add(GTK_CONTAINER(window),GTK_WIDGET(glarea));

	gdk_threads_add_idle(queue_render_idle,glarea);

	gtk_widget_show_all(GTK_WIDGET(window));
	return 1;
}

void windowGetSize(float* width, float* height)
{
	if(width == NULL || height == NULL)
	{
		return;
	}

	GtkAllocation* alloc = g_new(GtkAllocation, 1);
    gtk_widget_get_allocation(GTK_WIDGET(window), alloc);
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