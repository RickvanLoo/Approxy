library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

entity {{.EntityName}} is
generic (word_size: integer:={{.BitSize}}); 
Port ( 
A : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
B : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0));
end {{.EntityName}};

architecture Behavioral of {{.EntityName}} is
begin
	process(A,B) is
	begin
		case A is
            {{- range $rowindex, $row := .LUT}}
			when "{{$rowindex | indexconv}}" =>
				case B is
                    {{- range $itemindex, $val := $row}}
                    when "{{$itemindex | indexconv}}" =>
                        prod <= "{{$val | valconv}}";
                    {{- end}}
                    when others =>
                        prod <= (others => 'X');
                end case;
            {{- end}}
            when others =>
                prod <= (others => 'X');
		end case;
		end process;

end Behavioral;