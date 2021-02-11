rm -rf _tmp
mkdir _tmp
cd _tmp
echo ""

echo "[TEST] create new project:"
echo ""
go run ../runnr/runnr.go g new
echo ""
echo "-------------------------------------------------------"
echo ""

echo "[TEST] running help:"
echo ""
go run ../runnr/runnr.go -h
echo ""
echo "-------------------------------------------------------"
echo ""

echo "[TEST] running recompile and hello:"
echo ""
go run ../runnr/runnr.go hello -r
echo ""
echo "-------------------------------------------------------"
echo ""

echo "[TEST] test the example application:"
echo ""
cd ../
cd examples/example1
go run ../../runnr/runnr.go something -r
echo ""
echo "-------------------------------------------------------"
echo ""

echo "[TEST] can run it statically"
echo ""
cd ../../
go run ./runnr/runnr.go g s run ./_tmp/internal/cmd/runnr_local/main.go ./_tmp/statictest -- hello -r
echo ""
echo "-------------------------------------------------------"
echo ""
