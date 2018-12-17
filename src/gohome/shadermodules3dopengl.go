package gohome

import (
	"github.com/PucklaMotzer09/GLSLGenerator"
	"strings"
)

// 3D

const (
	ShaderVersion = "110"
)

var (
	Attributes3D = []glslgen.Variable{
		glslgen.Variable{"vec3", "highp", "vertex"},
		glslgen.Variable{"vec3", "highp", "normal"},
		glslgen.Variable{"vec2", "highp", "texCoord"},
		glslgen.Variable{"vec3", "highp", "tangent"},
	}

	AttributesInstanced3D = []glslgen.Variable{
		glslgen.Variable{"mat4", "highp", "transformMatrix3D"},
	}

	UniformModuleVertex3D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"mat4", "highp", "viewMatrix3D"},
			glslgen.Variable{"mat4", "highp", "inverseViewMatrix3D"},
			glslgen.Variable{"mat4", "highp", "projectionMatrix3D"},
		},
	}

	UniformNormalModuleVertex3D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"mat4", "highp", "transformMatrix3D"},
		},
	}

	CalculatePositionModule3D = glslgen.Module{
		Name: "calculatePosition",
		Body: "gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);",
	}

	SetOutputsModuleVertex3D = glslgen.Module{
		Name: "setOutputs",
		Body: `fragViewMatrix3D = viewMatrix3D;
			   fragInverseViewMatrix3D = inverseViewMatrix3D;`,
	}

	SetOutputsNormalModuleVertex3D = glslgen.Module{
		Name: "setOutputsNormal",
		Body: `fragPos =  (viewMatrix3D*transformMatrix3D*vec4(vertex,1.0)).xyz;
			   fragTexCoord = texCoord;
			   fragNormal =  (viewMatrix3D*transformMatrix3D*vec4(normal,0.0)).xyz;
			   vec3 norm = normalize(fragNormal);
			   vec3 tang = normalize((viewMatrix3D*transformMatrix3D*vec4(tangent,0.0)).xyz);
			   vec3 bitang = normalize(cross(norm,tang));
	
			   fragToTangentSpace = mat3(
				   tang.x,bitang.x,norm.x,
				   tang.y,bitang.y,norm.y,
				   tang.z,bitang.z,norm.z
			   );`,
	}

	SetOutputsNoUVModuleVertex3D = glslgen.Module{
		Name: "setOutputsNoUV",
		Body: `fragPos =  (transformMatrix3D*vec4(vertex,1.0)).xyz;
			   fragNormal =  (transformMatrix3D*vec4(normal,0.0)).xyz;`,
	}

	GlobalsFragment3D = []glslgen.Variable{
		glslgen.Variable{"float", "const highp", "shadowDistance = 50.0"},
		glslgen.Variable{"float", "const highp", "transitionDistance = 5.0"},
		glslgen.Variable{"float", "const highp", "bias = 0.005"},
		glslgen.Variable{"vec4", "highp", "finalDiffuseColor"},
		glslgen.Variable{"vec4", "highp", "finalSpecularColor"},
		glslgen.Variable{"vec4", "highp", "finalAmbientColor"},
		glslgen.Variable{"vec3", "highp", "norm"},
		glslgen.Variable{"vec3", "highp", "viewDir"},
	}

	InputsFragment3D = []glslgen.Variable{
		glslgen.Variable{"vec3", "highp", "fragPos"},
		glslgen.Variable{"vec3", "highp", "fragNormal"},
		glslgen.Variable{"mat4", "highp", "fragViewMatrix3D"},
		glslgen.Variable{"mat4", "highp", "fragInverseViewMatrix3D"},
	}

	InputsNormalFragment3D = []glslgen.Variable{
		glslgen.Variable{"vec2", "highp", "fragTexCoord"},
		glslgen.Variable{"mat3", "highp", "fragToTangentSpace"},
	}

	LightMakrosFragment3D = []glslgen.Makro{
		glslgen.Makro{"MAX_POINT_LIGHTS", "5"},
		glslgen.Makro{"MAX_DIRECTIONAL_LIGHTS", "2"},
		glslgen.Makro{"MAX_SPOT_LIGHTS", "1"},
		glslgen.Makro{"degToRad(deg)", "(deg/180.0*3.14159265359)"},
		glslgen.Makro{"MAX_SPECULAR_EXPONENT", "50.0"},
		glslgen.Makro{"MIN_SPECULAR_EXPONENT", "5.0"},
	}

	InitialiseModuleFragment3D = glslgen.Module{
		Name: "initialise",
		Body: `finalDiffuseColor = vec4(1.0,1.0,1.0,1.0);
			   finalSpecularColor = vec4(0.0);
			   finalAmbientColor = vec4(0.0);`,
	}

	InitialiseNormalModuleFragment3D = glslgen.Module{
		Name: "initialiseNormal",
		Body: `norm = normalize(fragToTangentSpace*fragNormal);
			   viewDir = normalize(fragToTangentSpace*(fragPos*-1.0));`,
	}

	InitialiseNoUVModuleFragment3D = glslgen.Module{
		Name: "initialiseNoUV",
		Body: `norm = normalize(fragNormal);
			   vec3 camPos = (fragInverseViewMatrix3D*vec4(0.0,0.0,0.0,1.0)).xyz;
			   viewDir = camPos - fragPos;`,
	}

	LightUniformsModule3D = glslgen.Module{
		Structs: []glslgen.Struct{
			glslgen.Struct{
				"Attentuation",
				[]glslgen.Variable{
					glslgen.Variable{"float", "highp", "constant"},
					glslgen.Variable{"float", "highp", "linear"},
					glslgen.Variable{"float", "highp", "quadratic"},
				},
			},
			glslgen.Struct{
				"PointLight",
				[]glslgen.Variable{
					glslgen.Variable{"vec3", "highp", "position"},
					glslgen.Variable{"vec3", "highp", "diffuseColor"},
					glslgen.Variable{"vec3", "highp", "specularColor"},
					glslgen.Variable{"Attentuation", "", "attentuation"},
				},
			},
			glslgen.Struct{
				"DirectionalLight",
				[]glslgen.Variable{
					glslgen.Variable{"vec3", "highp", "direction"},
					glslgen.Variable{"vec3", "highp", "diffuseColor"},
					glslgen.Variable{"vec3", "highp", "specularColor"},
					glslgen.Variable{"mat4", "highp", "lightSpaceMatrix"},
					glslgen.Variable{"bool", "", "castsShadows"},
					glslgen.Variable{"ivec2", "", "shadowMapSize"},
					glslgen.Variable{"float", "highp", "shadowDistance"},
				},
			},
			glslgen.Struct{
				"SpotLight",
				[]glslgen.Variable{
					glslgen.Variable{"vec3", "highp", "position"},
					glslgen.Variable{"vec3", "highp", "direction"},
					glslgen.Variable{"vec3", "highp", "diffuseColor"},
					glslgen.Variable{"vec3", "highp", "specularColor"},
					glslgen.Variable{"float", "highp", "innerCutOff"},
					glslgen.Variable{"float", "highp", "outerCutOff"},
					glslgen.Variable{"Attentuation", "", "attentuation"},
					glslgen.Variable{"mat4", "highp", "lightSpaceMatrix"},
					glslgen.Variable{"bool", "", "castsShadows"},
					glslgen.Variable{"ivec2", "", "shadowMapSize"},
				},
			},
		},

		Uniforms: []glslgen.Variable{
			glslgen.Variable{"int", "", "numPointLights"},
			glslgen.Variable{"int", "", "numDirectionalLights"},
			glslgen.Variable{"int", "", "numSpotLights"},
			glslgen.Variable{"vec3", "highp", "ambientLight"},
			glslgen.Variable{"PointLight", "", "pointLights[MAX_POINT_LIGHTS]"},
			glslgen.Variable{"DirectionalLight", "", "directionalLights[MAX_POINT_LIGHTS]"},
			glslgen.Variable{"sampler2D", "highp", "directionalLightsshadowmap[MAX_DIRECTIONAL_LIGHTS]"},
			glslgen.Variable{"SpotLight", "", "spotLights[MAX_SPOT_LIGHTS]"},
			glslgen.Variable{"sampler2D", "highp", "spotLightsshadowmap[MAX_SPOT_LIGHTS]"},
		},

		Functions: []glslgen.Function{
			glslgen.Function{
				"vec3 diffuseLighting(vec3 lightDir, vec3 diffuse)",
				`float diff = max(dot(norm,lightDir),0.0);
				 diffuse *= diff;
				 return diffuse;`,
			},
			glslgen.Function{
				"vec3 specularLighting(vec3 lightDir, vec3 specular)",
				`vec3 reflectDir = reflect(-lightDir, norm);
				 vec3 halfwayDir = normalize(lightDir + viewDir);
				 float spec = max(pow(max(dot(norm,halfwayDir),0.0),calculateShinyness(material.shinyness)),0.0);
				 specular *= spec;
				 return specular;`,
			},
			glslgen.Function{
				"void calculatePointLights()",
				` for (int i = 0;i<numPointLights&&i<MAX_POINT_LIGHTS;i++)
				 {
				 	calculatePointLight(pointLights[i],i);
				 }
			
				/*#if MAX_POINT_LIGHTS > 0
				if(numPointLights > 0)
					calculatePointLight(pointLights[0],0);
				#endif
				#if MAX_POINT_LIGHTS > 1
				if(numPointLights > 1)
					calculatePointLight(pointLights[1],1);
				#endif
				#if MAX_POINT_LIGHTS > 2
				if(numPointLights > 2)
					calculatePointLight(pointLights[2],2);
				#endif
				#if MAX_POINT_LIGHTS > 3
				if(numPointLights > 3)
					calculatePointLight(pointLights[3],3);
				#endif
				#if MAX_POINT_LIGHTS > 4
				if(numPointLights > 4)
					calculatePointLight(pointLights[4],4);
				#endif
				#if MAX_POINT_LIGHTS > 5
				if(numPointLights > 5)
					calculatePointLight(pointLights[5],5);
				#endif
				#if MAX_POINT_LIGHTS > 6
				if(numPointLights > 6)
					calculatePointLight(pointLights[6],6);
				#endif
				#if MAX_POINT_LIGHTS > 7
				if(numPointLights > 7)
					calculatePointLight(pointLights[7],7);
				#endif
				#if MAX_POINT_LIGHTS > 8
				if(numPointLights > 8)
					calculatePointLight(pointLights[8],8);
				#endif
				#if MAX_POINT_LIGHTS > 9
				if(numPointLights > 9)
					calculatePointLight(pointLights[9],9);
				#endif
				#if MAX_POINT_LIGHTS > 10
				if(numPointLights > 10)
					calculatePointLight(pointLights[10],10);
				#endif
				#if MAX_POINT_LIGHTS > 11
				if(numPointLights > 11)
					calculatePointLight(pointLights[11],11);
				#endif
				#if MAX_POINT_LIGHTS > 12
				if(numPointLights > 12)
					calculatePointLight(pointLights[12],12);
				#endif
				#if MAX_POINT_LIGHTS > 13
				if(numPointLights > 13)
					calculatePointLight(pointLights[13],13);
				#endif
				#if MAX_POINT_LIGHTS > 14
				if(numPointLights > 14)
					calculatePointLight(pointLights[14],14);
				#endif
				#if MAX_POINT_LIGHTS > 15
				if(numPointLights > 15)
					calculatePointLight(pointLights[15],15);
				#endif
				#if MAX_POINT_LIGHTS > 16
				if(numPointLights > 16)
					calculatePointLight(pointLights[16],16);
				#endif
				#if MAX_POINT_LIGHTS > 17
				if(numPointLights > 17)
					calculatePointLight(pointLights[17],17);
				#endif
				#if MAX_POINT_LIGHTS > 18
				if(numPointLights > 18)
					calculatePointLight(pointLights[18],18);
				#endif
				#if MAX_POINT_LIGHTS > 19
				if(numPointLights > 19)
					calculatePointLight(pointLights[19],19);
				#endif
				#if MAX_POINT_LIGHTS > 20
				if(numPointLights > 20)
					calculatePointLight(pointLights[20],20);
				#endif*/`,
			},
			glslgen.Function{
				"void calculateDirectionalLights()",
				`
				for (int i = 0;i<numDirectionalLights&&i<MAX_DIRECTIONAL_LIGHTS;i++)
				 {
					calculateDirectionalLight(directionalLights[i],i);
				 }
				/*#if MAX_DIRECTIONAL_LIGHTS > 0
				if(int(numDirectionalLights) > 0)
					calculateDirectionalLight(directionalLights[0],0);
				#endif
				#if MAX_DIRECTIONAL_LIGHTS > 1
				if(int(numDirectionalLights) > 1)
					calculateDirectionalLight(directionalLights[1],1);
				#endif
				#if MAX_DIRECTIONAL_LIGHTS > 2
				if(int(numDirectionalLights) > 2)
					calculateDirectionalLight(directionalLights[2],2);
				#endif
				#if MAX_DIRECTIONAL_LIGHTS > 3
				if(int(numDirectionalLights) > 3)
					calculateDirectionalLight(directionalLights[3],3);
				#endif
				#if MAX_DIRECTIONAL_LIGHTS > 4
				if(int(numDirectionalLights) > 4)
					calculateDirectionalLight(directionalLights[4],4);
				#endif
				#if MAX_DIRECTIONAL_LIGHTS > 5
				if(int(numDirectionalLights) > 5)
					calculateDirectionalLight(directionalLights[5],5);
				#endif
				#if MAX_DIRECTIONAL_LIGHTS > 6
				if(int(numDirectionalLights) > 6)
					calculateDirectionalLight(directionalLights[6],6);
				#endif
				#if MAX_DIRECTIONAL_LIGHTS > 7
				if(int(numDirectionalLights) > 7)
					calculateDirectionalLight(directionalLights[7],7);
				#endif
				#if MAX_DIRECTIONAL_LIGHTS > 8
				if(int(numDirectionalLights) > 8)
					calculateDirectionalLight(directionalLights[8],8);
				#endif*/`,
			},
			glslgen.Function{
				"void calculateSpotLights()",
				`for(int i=0; i<numSpotLights && i<MAX_SPOT_LIGHTS ; i++)
				 {
				 	calculateSpotLight(spotLights[i],i);
				 }
				/*#if MAX_SPOT_LIGHTS > 0
				if(int(numSpotLights) > 0)
					calculateSpotLight(spotLights[0],0);
				#endif
				#if MAX_SPOT_LIGHTS > 1
				if(int(numSpotLights) > 1)
					calculateSpotLight(spotLights[1],1);
				#endif
				#if MAX_SPOT_LIGHTS > 2
				if(int(numSpotLights) > 2)
					calculateSpotLight(spotLights[2],2);
				#endif
				#if MAX_SPOT_LIGHTS > 3
				if(int(numSpotLights) > 3)
					calculateSpotLight(spotLights[3],3);
				#endif
				#if MAX_SPOT_LIGHTS > 4
				if(int(numSpotLights) > 4)
					calculateSpotLight(spotLights[4],4);
				#endif
				#if MAX_SPOT_LIGHTS > 5
				if(int(numSpotLights) > 5)
					calculateSpotLight(spotLights[5],5);
				#endif
				#if MAX_SPOT_LIGHTS > 6
				if(int(numSpotLights) > 6)
					calculateSpotLight(spotLights[6],6);
				#endif
				#if MAX_SPOT_LIGHTS > 7
				if(int(numSpotLights) > 7)
					calculateSpotLight(spotLights[7],7);
				#endif
				#if MAX_SPOT_LIGHTS > 8
				if(int(numSpotLights) > 8)
					calculateSpotLight(spotLights[8],8);
				#endif
				#if MAX_SPOT_LIGHTS > 9
				if(int(numSpotLights) > 9)
					calculateSpotLight(spotLights[9],9);
				#endif
				#if MAX_SPOT_LIGHTS > 10
				if(int(numSpotLights) > 10)
					calculateSpotLight(spotLights[10],10);
				#endif
				#if MAX_SPOT_LIGHTS > 11
				if(int(numSpotLights) > 11)
					calculateSpotLight(spotLights[11],11);
				#endif
				#if MAX_SPOT_LIGHTS > 12
				if(int(numSpotLights) > 12)
					calculateSpotLight(spotLights[12],12);
				#endif
				#if MAX_SPOT_LIGHTS > 13
				if(int(numSpotLights) > 13)
					calculateSpotLight(spotLights[13],13);
				#endif
				#if MAX_SPOT_LIGHTS > 14
				if(int(numSpotLights) > 14)
					calculateSpotLight(spotLights[14],14);
				#endif
				#if MAX_SPOT_LIGHTS > 15
				if(int(numSpotLights) > 15)
					calculateSpotLight(spotLights[15],15);
				#endif
				#if MAX_SPOT_LIGHTS > 16
				if(int(numSpotLights) > 16)
					calculateSpotLight(spotLights[16],16);
				#endif*/`,
			},
			glslgen.Function{
				"void calculateAllLights()",
				`calculatePointLights();
				 calculateDirectionalLights();
				 calculateSpotLights();`,
			},
			glslgen.Function{
				"float calcAttentuation(vec3 lightPosition,Attentuation attentuation)",
				`float distance = distance(lightPosition,fragPos);
				float attent = 1.0/(attentuation.quadratic*distance*distance + attentuation.linear*distance + attentuation.constant);
				return attent;`,
			},
			glslgen.Function{
				"float calculateShinyness(float shinyness)",
				"return max(MAX_SPECULAR_EXPONENT*(pow(max(shinyness,0.0),-3.0)-1.0)+MIN_SPECULAR_EXPONENT,0.0);",
			},
		},
		Name: "calculateLights",
		Body: `finalDiffuseColor = vec4(0.0,0.0,0.0,1.0);
			   finalAmbientColor = vec4(ambientLight,1.0);
			   calculateAllLights();`,
	}

	LightCalcSpotAmountNormalModule3D = glslgen.Module{
		Functions: []glslgen.Function{
			glslgen.Function{
				"float calcSpotAmount(vec3 lightDir,vec3 lightDirection,SpotLight pl)",
				`float theta = dot(lightDir, normalize(fragToTangentSpace*lightDirection));
			float spotAmount = 0.0;
			float outerCutOff = cos(degToRad(pl.outerCutOff));
			float innerCutOff = cos(degToRad(pl.innerCutOff));
			float epsilon   = innerCutOff - outerCutOff;
			spotAmount = clamp((theta - outerCutOff) / epsilon, 0.0, 1.0);
		
			return spotAmount;`,
			},
		},
	}

	LightCalcSpotAmountNoUVModule3D = glslgen.Module{
		Functions: []glslgen.Function{
			glslgen.Function{
				"float calcSpotAmount(vec3 lightDir,vec3 lightDirection,SpotLight pl)",
				`float theta = dot(lightDir, lightDirection);
			float spotAmount = 0.0;
			float outerCutOff = cos(degToRad(pl.outerCutOff));
			float innerCutOff = cos(degToRad(pl.innerCutOff));
			float epsilon   = innerCutOff - outerCutOff;
			spotAmount = clamp((theta - outerCutOff) / epsilon, 0.0, 1.0);
		
			return spotAmount;`,
			},
		},
	}

	LightsAndShadowsFunctionsNoUV3D = glslgen.Module{
		Functions: []glslgen.Function{
			glslgen.Function{
				"float calcShadow(sampler2D shadowMap,mat4 lightSpaceMatrix,float shadowdistance,bool distanceTransition,ivec2 shadowMapSize)",
				`float distance = 0.0;
				if(distanceTransition)
				{
					distance = length(fragPos);
					distance = distance - (shadowdistance - transitionDistance);
					distance = distance / transitionDistance;
					distance = clamp(1.0-distance,0.0,1.0);
				}
				vec4 fragPosLightSpace = lightSpaceMatrix*vec4(fragPos,1.0);
				vec3 projCoords = clamp((fragPosLightSpace.xyz / fragPosLightSpace.w)*0.5+0.5,-1.0,1.0);
				float currentDepth = projCoords.z-bias;
				float shadowresult = 0.0;
				float closestDepth = texture2D(shadowMap, projCoords.xy).r;
				vec2 texelSize = 1.0 / vec2(shadowMapSize);
				for(int x = -1; x <= 1; ++x)
				{
					for(int y = -1; y <= 1; ++y)
					{
						float pcfDepth = texture2D(shadowMap, projCoords.xy + vec2(x, y) * texelSize).r; 
						shadowresult += currentDepth > pcfDepth ? 0.0 : 1.0;        
					}    
				}
				shadowresult /= 9.0;
				if(distanceTransition)
				{
					shadowresult = 1.0 - (1.0-shadowresult)*distance;
				}
				return shadowresult;`,
			},
		},
		Name: "lightsAndShadowCalculationNoUV",
	}

	LightsAndShadowsFunctions3D = glslgen.Module{
		Functions: []glslgen.Function{
			glslgen.Function{
				"float calcShadow(sampler2D shadowMap,mat4 lightSpaceMatrix,float shadowdistance,bool distanceTransition,ivec2 shadowMapSize)",
				`float distance = 0.0;
				if(distanceTransition)
				{
					distance = length(fragPos);
					distance = distance - (shadowdistance - transitionDistance);
					distance = distance / transitionDistance;
					distance = clamp(1.0-distance,0.0,1.0);
				}
				vec4 fragPosLightSpace = lightSpaceMatrix*fragInverseViewMatrix3D*vec4(fragPos,1.0);
				vec3 projCoords = clamp((fragPosLightSpace.xyz / fragPosLightSpace.w)*0.5+0.5,-1.0,1.0);
				float currentDepth = projCoords.z-bias;
				float shadowresult = 0.0;
				float closestDepth = texture2D(shadowMap, projCoords.xy).r;
				vec2 texelSize = 1.0 / vec2(shadowMapSize);
				for(int x = -1; x <= 1; ++x)
				{
					for(int y = -1; y <= 1; ++y)
					{
						float pcfDepth = texture2D(shadowMap, projCoords.xy + vec2(x, y) * texelSize).r; 
						shadowresult += currentDepth > pcfDepth ? 0.0 : 1.0;        
					}    
				}
				shadowresult /= 9.0;
				if(distanceTransition)
				{
					shadowresult = 1.0 - (1.0-shadowresult)*distance;
				}
				return shadowresult;`,
			},
		},
		Name: "lightsAndShadowCalculation",
	}

	calcPointLightFunc = glslgen.Function{
		"void calculatePointLight(PointLight pl,int index)",
		`vec3 lightPosition = (fragViewMatrix3D*vec4(pl.position,1.0)).xyz;
				vec3 lightDir = normalize(fragToTangentSpace*(lightPosition - fragPos));
			
			
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);
			
				// Specular
				vec3 specular = specularLighting(lightDir,pl.specularColor);
			
				// Attentuation
				float attent = calcAttentuation(lightPosition,pl.attentuation);
			
				diffuse *= attent;
				specular *= attent;
			
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
	}

	LightsAndShadowsCalculationModule3D = glslgen.Module{
		Functions: []glslgen.Function{
			calcPointLightFunc,
			glslgen.Function{
				"void calculateDirectionalLight(DirectionalLight dl,int index)",
				`vec3 lightDirection = (fragViewMatrix3D*vec4(dl.direction*-1.0,0.0)).xyz;
				vec3 lightDir = normalize(fragToTangentSpace*lightDirection);
				
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,dl.diffuseColor);
				
				// Specular
				vec3 specular = specularLighting(lightDir,dl.specularColor);
				
				// Shadow
				float shadow = dl.castsShadows ? calcShadow(directionalLightsshadowmap[index],dl.lightSpaceMatrix,dl.shadowDistance,true,dl.shadowMapSize) : 1.0;
				
				diffuse *= shadow;
				specular *= shadow;
			
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
			},
			glslgen.Function{
				"void calculateSpotLight(SpotLight pl,int index)",
				`vec3 lightPosition = (fragViewMatrix3D*vec4(pl.position,1.0)).xyz;
				vec3 lightDirection = (fragViewMatrix3D*vec4(pl.direction*-1.0,0.0)).xyz;
				vec3 lightDir = normalize(fragToTangentSpace*(lightPosition-fragPos));
			
				// Spotamount
				float spotAmount = calcSpotAmount(lightDir,lightDirection,pl);
			
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);
			
				// Specular
				vec3 specular = specularLighting(lightDir,pl.specularColor);
			
				// Attentuation
				float attent = calcAttentuation(lightPosition,pl.attentuation);
			
				// Shadow
				float shadow = pl.castsShadows ? calcShadow(spotLightsshadowmap[index],pl.lightSpaceMatrix,50.0,false,pl.shadowMapSize) : 1.0;
				// float shadow = 1.0;
			
				diffuse *= attent * spotAmount * shadow;
				specular *= attent * spotAmount * shadow;
			
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
			},
		},
	}

	LightCalculationModel3D = glslgen.Module{
		Functions: []glslgen.Function{
			calcPointLightFunc,
			glslgen.Function{
				"void calculateDirectionalLight(DirectionalLight dl,int index)",
				`vec3 lightDirection = (fragViewMatrix3D*vec4(dl.direction*-1.0,0.0)).xyz;
				vec3 lightDir = normalize(fragToTangentSpace*lightDirection);
				
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,dl.diffuseColor);
				
				// Specular
				vec3 specular = specularLighting(lightDir,dl.specularColor);
				
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
			},
			glslgen.Function{
				"void calculateSpotLight(SpotLight pl,int index)",
				`vec3 lightPosition = (fragViewMatrix3D*vec4(pl.position,1.0)).xyz;
				vec3 lightDirection = (fragViewMatrix3D*vec4(pl.direction*-1.0,0.0)).xyz;
				vec3 lightDir = normalize(fragToTangentSpace*(lightPosition-fragPos));
			
				// Spotamount
				float spotAmount = calcSpotAmount(lightDir,lightDirection,pl);
			
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);
			
				// Specular
				vec3 specular = specularLighting(lightDir,pl.specularColor);
			
				// Attentuation
				float attent = calcAttentuation(lightPosition,pl.attentuation);
			
				diffuse *= attent * spotAmount;
				specular *= attent * spotAmount;
			
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
			},
		},
	}

	calcPointLightNoUVFunc = glslgen.Function{
		"void calculatePointLight(PointLight pl,int index)",
		`vec3 lightPosition = pl.position;
		vec3 lightDir = normalize(lightPosition - fragPos);
	
	
		// Diffuse
		vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);
	
		// Specular
		vec3 specular = specularLighting(lightDir,pl.specularColor);
	
		// Attentuation
		float attent = calcAttentuation(lightPosition,pl.attentuation);
	
		diffuse *= attent;
		specular *= attent;
	
		finalDiffuseColor += vec4(diffuse,0.0);
		finalSpecularColor += vec4(specular,0.0);`,
	}

	LightsAndShadowsCalculationNoUVModule3D = glslgen.Module{
		Functions: []glslgen.Function{
			calcPointLightNoUVFunc,
			glslgen.Function{
				"void calculateDirectionalLight(DirectionalLight dl,int index)",
				`vec3 lightDirection = -dl.direction;
				vec3 lightDir = normalize(lightDirection);
				
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,dl.diffuseColor);
				
				// Specular
				vec3 specular = specularLighting(lightDir,dl.specularColor);
				
				// Shadow
				float shadow = dl.castsShadows ? calcShadow(directionalLightsshadowmap[index],dl.lightSpaceMatrix,dl.shadowDistance,true,dl.shadowMapSize) : 1.0;
				
				diffuse *= shadow;
				specular *= shadow;
			
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
			},
			glslgen.Function{
				"void calculateSpotLight(SpotLight pl,int index)",
				`vec3 lightPosition = pl.position;
				vec3 lightDirection = -pl.direction;
				vec3 lightDir = normalize(lightPosition-fragPos);
			
				// Spotamount
				float spotAmount = calcSpotAmount(lightDir,lightDirection,pl);
			
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);
			
				// Specular
				vec3 specular = specularLighting(lightDir,pl.specularColor);
			
				// Attentuation
				float attent = calcAttentuation(lightPosition,pl.attentuation);
			
				// Shadow
				float shadow = pl.castsShadows ? calcShadow(spotLightsshadowmap[index],pl.lightSpaceMatrix,50.0,false,pl.shadowMapSize) : 1.0;
				// float shadow = 1.0;
			
				diffuse *= attent * spotAmount * shadow;
				specular *= attent * spotAmount * shadow;
			
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
			},
		},
	}

	LightCalculationNoUVModule3D = glslgen.Module{
		Functions: []glslgen.Function{
			calcPointLightNoUVFunc,
			glslgen.Function{
				"void calculateDirectionalLight(DirectionalLight dl,int index)",
				`vec3 lightDirection = -dl.direction;
				vec3 lightDir = normalize(lightDirection);
				
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,dl.diffuseColor);
				
				// Specular
				vec3 specular = specularLighting(lightDir,dl.specularColor);
				
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
			},
			glslgen.Function{
				"void calculateSpotLight(SpotLight pl,int index)",
				`vec3 lightPosition = pl.position;
				vec3 lightDirection = -pl.direction;
				vec3 lightDir = normalize(lightPosition-fragPos);
			
				// Spotamount
				float spotAmount = calcSpotAmount(lightDir,lightDirection,pl);
			
				// Diffuse
				vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);
			
				// Specular
				vec3 specular = specularLighting(lightDir,pl.specularColor);
			
				// Attentuation
				float attent = calcAttentuation(lightPosition,pl.attentuation);
			
				diffuse *= attent * spotAmount;
				specular *= attent * spotAmount;
			
				finalDiffuseColor += vec4(diffuse,0.0);
				finalSpecularColor += vec4(specular,0.0);`,
			},
		},
	}

	MaterialModule3D = glslgen.Module{
		Structs: []glslgen.Struct{
			glslgen.Struct{
				Name: "Material",
				Variables: []glslgen.Variable{
					glslgen.Variable{"vec3", "highp", "diffuseColor"},
					glslgen.Variable{"vec3", "highp", "specularColor"},
					glslgen.Variable{"bool", "", "DiffuseTextureLoaded"},
					glslgen.Variable{"bool", "", "SpecularTextureLoaded"},
					glslgen.Variable{"bool", "", "NormalMapLoaded"},
					glslgen.Variable{"float", "highp", "shinyness"},
					glslgen.Variable{"float", "highp", "transparency"},
				},
			},
		},
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"Material", "", "material"},
		},
		Name: "materialCalculation",
		Body: `finalDiffuseColor *= vec4(material.diffuseColor,material.transparency);
			   finalSpecularColor *= vec4(material.specularColor,0.0);
			   finalAmbientColor *= vec4(material.diffuseColor,0.0);`,
	}

	DiffuseTextureModule3D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"sampler2D", "highp", "materialdiffuseTexture"},
		},
		Name: "diffuseTextureCalculation",
		Body: `vec4 texDifCol;
			   if(material.DiffuseTextureLoaded)
				  texDifCol = texture2D(materialdiffuseTexture,fragTexCoord);
			   else
				  texDifCol = vec4(1.0);
			   finalDiffuseColor *= texDifCol;
			   finalAmbientColor *= texDifCol;`,
	}

	SpecularTextureModule3D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"sampler2D", "highp", "materialspecularTexture"},
		},
		Name: "specularTextureCalculation",
		Body: `vec4 texSpecCol;
			   if(material.SpecularTextureLoaded)
			     texSpecCol = texture2D(materialspecularTexture,fragTexCoord);
			   else
				 texSpecCol = vec4(1.0);
			   finalSpecularColor *= texSpecCol;`,
	}

	NormalMapModule3D = glslgen.Module{
		Uniforms: []glslgen.Variable{
			glslgen.Variable{"sampler2D", "highp", "materialnormalMap"},
		},
		Name: "normalMapCalculation",
		Body: `if(material.NormalMapLoaded)
				  norm = normalize(2.0*(texture2D(materialnormalMap,fragTexCoord)).xyz-1.0);`,
	}

	FinalModuleFragment3D = glslgen.Module{
		Name: "finalCalculation",
		Body: `if(finalDiffuseColor.a < 0.1)
				  discard;
			   gl_FragColor = finalDiffuseColor + finalSpecularColor + finalAmbientColor;`,
	}
)

func LoadGeneratedShader3D(shader_type uint8, flags uint32) Shader {
	n, v, f := GenerateShader3D(shader_type, flags)
	return ls(n, v, f)
}

func GenerateShaderSource3D() {
	n, v, f := GenerateShader3D(SHADER_TYPE_3D, 0)
	ls(n, v, f)
	n, v, f = GenerateShader3D(SHADER_TYPE_3D, SHADER_FLAG_NOUV)
	ls(n, v, f)
	n, v, f = GenerateShader3D(SHADER_TYPE_3D, SHADER_FLAG_INSTANCED)
	ls(n, v, f)
	n, v, f = GenerateShader3D(SHADER_TYPE_3D, SHADER_FLAG_INSTANCED|SHADER_FLAG_NOUV)
	ls(n, v, f)
}

func ls(n, v, f string) Shader {
	return ResourceMgr.LoadShaderSource(n, v, f, "", "", "", "")
}

const (
	SHADER_FLAG_INSTANCED   uint32 = (1 << 0)
	SHADER_FLAG_NOUV        uint32 = (1 << 1)
	SHADER_FLAG_NO_SHADOWS  uint32 = (1 << 2)
	SHADER_FLAG_NO_LIGHTING uint32 = (1 << 3)
	SHADER_FLAG_NO_DIFTEX   uint32 = (1 << 4)
	SHADER_FLAG_NO_SPECTEX  uint32 = (1 << 5)
	SHADER_FLAG_NO_NORMAP   uint32 = (1 << 6)
	NUM_FLAGS_3D            uint32 = 7

	SHADER_FLAG_NO_KEYCOLOR       uint32 = (1 << 1)
	SHADER_FLAG_NO_MODCOLOR       uint32 = (1 << 2)
	SHADER_FLAG_NO_FLIP           uint32 = (1 << 3)
	SHADER_FLAG_NO_TEXTURE_REGION uint32 = (1 << 4)
	SHADER_FLAG_NO_DEPTH          uint32 = (1 << 5)
	SHADER_FLAG_NO_TEXTURE        uint32 = (1 << 6)
	SHADER_FLAG_DEPTHMAP          uint32 = (1 << 7)
	NUM_FLAGS_2D                  uint32 = 8
)

var (
	FLAG_NAMES_3D = [NUM_FLAGS_3D]string{
		"Instanced",
		"NoUV",
		"NoShadows",
		"NoShadows NoLighting",
		"NoDiftex",
		"NoSpectex",
		"NoNormap",
	}

	FLAG_NAMES_2D = [NUM_FLAGS_2D]string{
		"Instanced",
		"NoKeyColor",
		"NoModColor",
		"NoFlip",
		"NoTextureRegion",
		"NoDepth",
		"NoTexture",
		"DepthMap",
	}
)

func GetShaderName3D(flags uint32) string {
	var n = "3D"
	for i := uint32(0); i < NUM_FLAGS_3D; i++ {
		if flags&(1<<i) != 0 {
			n += " " + FLAG_NAMES_3D[i]
		}
	}
	return n
}

func GenerateShader3D(shader_type uint8, flags uint32) (n, v, f string) {
	if shader_type == SHADER_TYPE_3D {
		startFlags := flags
		if !Render.HasFunctionAvailable("INSTANCED") {
			flags &= ^SHADER_FLAG_INSTANCED
		}
		if flags&SHADER_FLAG_NO_LIGHTING != 0 {
			flags |= SHADER_FLAG_NO_SHADOWS
		}
		if flags&SHADER_FLAG_NOUV != 0 {
			flags |= SHADER_FLAG_NO_DIFTEX | SHADER_FLAG_NO_SPECTEX | SHADER_FLAG_NO_NORMAP
		}

		rname := Render.GetName()
		if rname == "OpenGLES2" {
			flags |= SHADER_FLAG_NO_SHADOWS | SHADER_FLAG_NO_NORMAP
		}

		var vertex glslgen.VertexGenerator
		var fragment glslgen.FragmentGenerator

		if strings.Contains(rname, "OpenGLES") {
			vertex.SetVersion("100")
			fragment.SetVersion("100")
		} else {
			vertex.SetVersion(ShaderVersion)
			fragment.SetVersion(ShaderVersion)
		}

		vertex.AddAttributes(Attributes3D)
		if flags&SHADER_FLAG_INSTANCED != 0 {
			vertex.AddAttributes(AttributesInstanced3D)
		}
		vertex.AddOutputs(InputsFragment3D)
		if flags&SHADER_FLAG_NOUV == 0 && rname != "OpenGLES2" {
			vertex.AddOutputs(InputsNormalFragment3D)
		}
		if rname == "OpenGLES2" {
			vertex.AddOutput(glslgen.Variable{"vec2", "highp", "fragTexCoord"})
		}
		vertex.AddModule(UniformModuleVertex3D)
		if flags&SHADER_FLAG_INSTANCED == 0 {
			vertex.AddModule(UniformNormalModuleVertex3D)
		}
		vertex.AddModule(CalculatePositionModule3D)
		vertex.AddModule(SetOutputsModuleVertex3D)
		if flags&SHADER_FLAG_NOUV == 0 && rname != "OpenGLES2" {
			vertex.AddModule(SetOutputsNormalModuleVertex3D)
		} else {
			vertex.AddModule(SetOutputsNoUVModuleVertex3D)
		}

		if flags&SHADER_FLAG_NO_LIGHTING == 0 {
			fragment.AddMakros(LightMakrosFragment3D)
		}
		fragment.AddGlobals(GlobalsFragment3D)

		fragment.AddInputs(InputsFragment3D)
		if flags&SHADER_FLAG_NOUV == 0 && rname != "OpenGLES2" {
			fragment.AddInputs(InputsNormalFragment3D)
		}
		if flags&SHADER_FLAG_NOUV == 0 && rname == "OpenGLES2" {
			fragment.AddInput(glslgen.Variable{"vec2", "highp", "fragTexCoord"})
		}
		fragment.AddModule(InitialiseModuleFragment3D)
		if flags&SHADER_FLAG_NOUV == 0 && rname != "OpenGLES2" {
			fragment.AddModule(InitialiseNormalModuleFragment3D)
		} else {
			fragment.AddModule(InitialiseNoUVModuleFragment3D)
		}
		if flags&(SHADER_FLAG_NO_NORMAP|SHADER_FLAG_NOUV) == 0 {
			fragment.AddModule(NormalMapModule3D)
		}
		if flags&SHADER_FLAG_NO_LIGHTING == 0 {
			fragment.AddModule(LightUniformsModule3D)
			if flags&SHADER_FLAG_NOUV == 0 && rname != "OpenGLES2" {
				fragment.AddModule(LightCalcSpotAmountNormalModule3D)
			} else {
				fragment.AddModule(LightCalcSpotAmountNoUVModule3D)
			}
			if flags&SHADER_FLAG_NO_SHADOWS == 0 {
				if flags&SHADER_FLAG_NOUV == 0 && rname != "OpenGLES2" {
					fragment.AddModule(LightsAndShadowsFunctions3D)
					fragment.AddModule(LightsAndShadowsCalculationModule3D)
				} else {
					fragment.AddModule(LightsAndShadowsFunctionsNoUV3D)
					fragment.AddModule(LightsAndShadowsCalculationNoUVModule3D)
				}
			} else {
				if flags&SHADER_FLAG_NOUV == 0 && rname != "OpenGLES2" {
					fragment.AddModule(LightCalculationModel3D)
				} else {
					fragment.AddModule(LightCalculationNoUVModule3D)
				}
			}
		}
		fragment.AddModule(MaterialModule3D)
		if flags&SHADER_FLAG_NO_DIFTEX == 0 {
			fragment.AddModule(DiffuseTextureModule3D)
		}
		if flags&SHADER_FLAG_NO_SPECTEX == 0 {
			fragment.AddModule(SpecularTextureModule3D)
		}
		fragment.AddModule(FinalModuleFragment3D)

		v = vertex.String()
		f = fragment.String()
		n = GetShaderName3D(startFlags)
	} else {
		n, v, f = generateShaderShape3D()
	}

	return
}
