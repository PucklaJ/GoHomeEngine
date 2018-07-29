#include "includes.h"
#include "gtkccallbacksheader.h"
#include "_cgo_export.h"

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
	Window = (GtkWindow*)gtk_window_new(GTK_WINDOW_TOPLEVEL);
	gtk_widget_set_size_request(GTK_WIDGET(Window),1,1);
	gtk_window_resize(Window,width,height);
	gtk_window_set_title(Window,title);
	gtk_widget_set_events(GTK_WIDGET(Window), GDK_POINTER_MOTION_MASK|GDK_SCROLL_MASK);

	g_signal_connect(GTK_WIDGET(Window),"delete-event",G_CALLBACK(gtkgo_quit_c),NULL);
	
	g_signal_connect(GTK_WIDGET(Window),"key-press-event",G_CALLBACK(gtkgo_gl_area_key_press_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"key-release-event",G_CALLBACK(gtkgo_gl_area_key_release_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"button-press-event",G_CALLBACK(gtkgo_gl_area_button_press_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"button-release-event",G_CALLBACK(gtkgo_gl_area_button_release_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"motion-notify-event",G_CALLBACK(gtkgo_gl_area_motion_notify_c),NULL);
	g_signal_connect(GTK_WIDGET(Window),"scroll-event",G_CALLBACK(gtkgo_gl_area_scroll_c),NULL);

    createGLArea();

	gtk_widget_show_all(GTK_WIDGET(Window));
	return 1;
}

void createGLArea()
{
	GLarea = (GtkGLArea*)gtk_gl_area_new();
	g_signal_connect(GTK_WIDGET(GLarea),"render",G_CALLBACK(gtkgo_gl_area_render_c),NULL);
	g_signal_connect(GTK_WIDGET(GLarea),"realize",G_CALLBACK(gtkgo_gl_area_realize_c),NULL);

	gtk_gl_area_set_has_depth_buffer(GLarea,TRUE);

	gdk_threads_add_idle(queue_render_idle,GLarea);
}

void windowSetSize(float width, float height)
{
	gtk_window_resize(Window,width,height);
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

void windowHideCursor()
{
	gdk_window_set_cursor(gtk_widget_get_window(GTK_WIDGET(Window)),gdk_cursor_new_from_name(gdk_display_get_default(),"none"));
}

void windowDisableCursor()
{
	// gdk_window_set_cursor(gtk_widget_get_window(GTK_WIDGET(Window)),gdk_cursor_new_from_name(gdk_display_get_default(),"none"));
	if(gdk_seat_grab(
				gdk_display_get_default_seat(gdk_display_get_default()),
				gtk_widget_get_window(GTK_WIDGET(Window)),
				GDK_SEAT_CAPABILITY_ALL_POINTING,
				FALSE,
				gdk_cursor_new_from_name(gdk_display_get_default(),"none"),
				NULL,
				NULL,
				NULL
	) != GDK_GRAB_SUCCESS)
	{
		g_print("Error disabling cursor\n");
	}
}

void windowShowCursor()
{
	gdk_seat_ungrab(gdk_display_get_default_seat(gdk_display_get_default()));
	gdk_window_set_cursor(gtk_widget_get_window(GTK_WIDGET(Window)),gdk_cursor_new_for_display(gdk_display_get_default(),GDK_ARROW));
}

int windowCursorShown()
{
	GdkCursor* cursor = gdk_window_get_cursor(gtk_widget_get_window(GTK_WIDGET(Window)));
	g_print("WindowCursorShown");
	if(cursor == NULL)
		g_print("Cursor is null");
	else if(gdk_cursor_get_cursor_type(cursor) == GDK_BLANK_CURSOR)
		g_print("Cursor is blank");
	return cursor != NULL && gdk_cursor_get_cursor_type(cursor) != GDK_BLANK_CURSOR;
}
int windowCursorHidden()
{
	GdkCursor* cursor = gdk_window_get_cursor(gtk_widget_get_window(GTK_WIDGET(Window)));
	g_print("WindowCursorHidden");
	if(cursor != NULL)
		g_print("Cursor is not null");
	else if(gdk_cursor_get_cursor_type(cursor) != GDK_BLANK_CURSOR)
		g_print("Cursor is not blank");
	return cursor == NULL || gdk_cursor_get_cursor_type(cursor) == GDK_BLANK_CURSOR;
}
int windowCursorDisabled()
{
	return windowCursorHidden();
}

void addGLAreaToWindow()
{
	addGLAreaToContainer(GTK_CONTAINER(Window));
}

void addGLAreaToContainer(GtkContainer* container)
{
	gtk_container_add(container,GTK_WIDGET(GLarea));
	gtk_widget_show(GTK_WIDGET(GLarea));
}

GtkContainer* boxToContainer(GtkBox* box)
{
	return GTK_CONTAINER(box);
}

GtkContainer* windowToContainer(GtkWindow* window)
{
	return GTK_CONTAINER(window);
}

GtkWidget* boxToWidget(GtkBox* box)
{
	return GTK_WIDGET(box);
}

GtkBox* widgetToBox(GtkWidget* widget)
{
	return GTK_BOX(widget);
}

GtkWidget* glareaToWidget(GtkGLArea* area)
{
	return GTK_WIDGET(area);
}

GtkWidget* buttonToWidget(GtkButton* button)
{
    return GTK_WIDGET(button);
}

GtkContainer* buttonToContainer(GtkButton* button)
{
    return GTK_CONTAINER(button);
}

GtkButton* widgetToButton(GtkWidget* widget)
{
    return GTK_BUTTON(widget);
}

void signalConnectButton(GtkButton* button,char* signal, int id)
{
    ButtonSignalUserData* bsud = (ButtonSignalUserData*)malloc(sizeof(ButtonSignalUserData));
    bsud->id = id;
    bsud->signal = signal;
    g_signal_connect(GTK_WIDGET(button),signal,G_CALLBACK(gtkgo_button_signal_c),bsud);
}