// Info de los modulos
#include <linux/module.h>
// Info del kernel en tiempo real
#include <linux/kernel.h>
//para la inf de la ram
#include <linux/mm.h>


// Headers para modulos
#include <linux/init.h>
// Header necesario para proc_fs
#include <linux/proc_fs.h>
// Para dar acceso al usuario
#include <asm/uaccess.h>
// Para manejar el directorio /proc
#include <linux/seq_file.h>

const long minute = 60;
const long hours = minute * 60;
const long day = hours * 24;
const long megabyte = 1024 * 1024;

//obtiene estadisticas del sistema
struct sysinfo s1;

static void init_meminfo(void)
{
    si_meminfo(&s1);
}


MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo de RAM para el Lab de Sopes 1");
MODULE_AUTHOR("Pablo Batz");

static int escribir_archivo(struct seq_file *archivo, void *v) {
    init_meminfo();
    unsigned long totalram = (s1.totalram*4);
    unsigned long freeram = (s1.freeram*4);
    unsigned long usedram = totalram - freeram;
    unsigned long porcentaje = (usedram*100)/totalram;

    seq_printf(archivo, "{\n");
    seq_printf(archivo, "\"total_memory\": %lu,\n" ,totalram/1024);
    seq_printf(archivo, "\"free_memory\": %lu,\n", freeram/1024);
    seq_printf(archivo, "\"porcentaje\": %lu,\n", porcentaje);
    seq_printf(archivo, "\"used_memory\": %lu\n", ((totalram - freeram)*100)/totalram);
    seq_printf(archivo, "}\n");
   
    return 0;
}

//Funcion que se ejecutara cada vez que se lea el archivo con el comando CAT
static int al_abrir(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_archivo, NULL);
}

//Si el kernel es 5.6 o mayor se usa la estructura proc_ops
static struct proc_ops operaciones =
{
    .proc_open = al_abrir,
    .proc_read = seq_read
};

//Funcion a ejecuta al insertar el modulo en el kernel con insmod
static int _insert(void)
{
    proc_create("ram_so1_1s2024", 0, NULL, &operaciones);
    printk(KERN_INFO "Pablo Javier Batz Contreras\n");
    return 0;
}

//Funcion a ejecuta al remover el modulo del kernel con rmmod
static void _remove(void)
{
    remove_proc_entry("ram_so1_1s2024", NULL);
    printk(KERN_INFO "Segundo Semestre 2024\n");
}

module_init(_insert);
module_exit(_remove);