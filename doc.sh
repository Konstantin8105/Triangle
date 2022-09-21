#!/bin/bash
rm -f *.out
rm -f *.test

echo -e "# triangle"
echo -e "Go interface of triangulation C Triangle"
echo -e "\`\`\`\n"

go doc -all .

echo -e "\`\`\`"
