curl http://127.0.0.1:8080/debug/pprof/heap -o mem_out.txt
curl http://127.0.0.1:8080/debug/pprof/profile?seconds=5 -o cpu_out.txt

go tool pprof -svg -alloc_objects pprof_1.exe mem_out.txt > mem_ao.svg
go tool pprof -svg pprof_1.exe cpu_out.txt > cpu.svg
